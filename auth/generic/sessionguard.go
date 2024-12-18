package generic

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/wolftotem4/golava-core/auth"
	"github.com/wolftotem4/golava-core/cookie"
	"github.com/wolftotem4/golava-core/session"
	"github.com/wolftotem4/golava-core/util"
)

type SessionGuard struct {
	Name             string
	Session          *session.SessionManager
	Cookie           cookie.IEncryptableCookieManager
	Request          *http.Request
	RememberDuration time.Duration
	Provider         auth.UserProvider

	user auth.Authenticatable
}

func (sg *SessionGuard) User() auth.Authenticatable {
	return sg.user
}

func (sg *SessionGuard) GetName() string {
	return fmt.Sprintf("login_%s", sg.Name)
}

func (sg *SessionGuard) GetRecallerName() string {
	return fmt.Sprintf("remember_%s", sg.Name)
}

func (sg *SessionGuard) SetUser(user auth.Authenticatable) {
	sg.user = user
}

func (sg *SessionGuard) ID() any {
	if sg.user == nil {
		return nil
	}

	return sg.user.GetAuthIdentifier()
}

func (sg *SessionGuard) Validate(ctx context.Context, credentials map[string]any) (bool, error) {
	_, valid, err := sg.validate(ctx, credentials)
	return valid, err
}

func (sg *SessionGuard) validate(ctx context.Context, credentials map[string]any) (auth.Authenticatable, bool, error) {
	user, err := sg.Provider.RetrieveByCredentials(ctx, credentials)
	if errors.Is(err, auth.ErrUserNotFound) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	valid, err := sg.Provider.ValidateCredentials(ctx, user, credentials)
	return user, valid, err
}

func (sg *SessionGuard) Check() bool {
	return sg.user != nil
}

func (sg *SessionGuard) HasUser() bool {
	return sg.user != nil
}

func (sg *SessionGuard) Guest() bool {
	return !sg.Check()
}

func (sg *SessionGuard) Attempt(ctx context.Context, credentials map[string]any, remember bool) (bool, error) {
	user, valid, err := sg.validate(ctx, credentials)
	if err != nil {
		return false, err
	} else if !valid {
		return false, nil
	}

	err = sg.Provider.RehashPasswordIfRequired(ctx, user, credentials, false)
	if err != nil {
		return false, err
	}

	err = sg.Login(ctx, user, remember)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Log a user into the application without sessions or cookies.
func (sg *SessionGuard) Once(ctx context.Context, credentials map[string]any) (bool, error) {
	user, valid, err := sg.validate(ctx, credentials)
	if err != nil {
		return false, err
	} else if !valid {
		return false, nil
	}

	err = sg.Provider.RehashPasswordIfRequired(ctx, user, credentials, false)
	if err != nil {
		return false, err
	}

	sg.SetUser(user)
	return true, nil
}

// Log the given user ID into the application without sessions or cookies.
func (sg *SessionGuard) OnceUsingID(ctx context.Context, id any) (bool, error) {
	user, err := sg.Provider.RetrieveById(ctx, id)
	if err != nil {
		return false, err
	}

	sg.SetUser(user)
	return true, nil
}

func (sg *SessionGuard) Login(ctx context.Context, user auth.Authenticatable, remember bool) error {
	err := sg.updateSession(ctx, user.GetAuthIdentifier())
	if err != nil {
		return err
	}

	if remember {
		err = sg.ensureRememberTokenIsSet(ctx, user)
		if err != nil {
			return err
		}
		sg.createRecaller(user)
	}

	sg.Cookie.Encryption().Set(
		sg.Session.Name,
		sg.Session.Store.ID,
		cookie.WithMaxAge(int(sg.Session.Lifetime.Seconds())),
	)

	sg.SetUser(user)

	return nil
}

func (sg *SessionGuard) ensureRememberTokenIsSet(ctx context.Context, user auth.Authenticatable) error {
	if user.GetRememberToken() == "" {
		return sg.cycleRememberToken(ctx, user)
	}
	return nil
}

func (sg *SessionGuard) cycleRememberToken(ctx context.Context, user auth.Authenticatable) error {
	token := util.RandomToken(45)
	user.SetRememberToken(token)
	return sg.Provider.UpdateRememberToken(ctx, user, token)
}

func (sg *SessionGuard) createRecaller(user auth.Authenticatable) {
	recaller := auth.NewRecallerString(
		user.GetAuthIdentifier(),
		user.GetRememberToken(),
		user.GetAuthPassword(),
	)
	sg.Cookie.Encryption().Set(
		sg.GetRecallerName(),
		recaller,
		cookie.WithMaxAge(int(sg.RememberDuration.Seconds())),
	)
}

func (sg *SessionGuard) updateSession(ctx context.Context, id any) error {
	sg.Session.Store.Put(sg.GetName(), id)
	return sg.Session.Store.Migrate(ctx, true)
}

func (sg *SessionGuard) LoginUsingID(ctx context.Context, id any, remember bool) error {
	user, err := sg.Provider.RetrieveById(ctx, id)
	if err != nil {
		return err
	}
	return sg.Login(ctx, user, remember)
}

func (sg *SessionGuard) Logout(ctx context.Context) error {
	user := sg.User()

	sg.clearUserDataFromStorage()

	if user != nil && user.GetRememberToken() != "" {
		err := sg.cycleRememberToken(ctx, user)
		if err != nil {
			return err
		}
	}

	sg.SetUser(nil)

	return nil
}

func (sg *SessionGuard) GetRecaller() (auth.Recaller, error) {
	recaller, err := sg.Cookie.Encryption().Get(sg.GetRecallerName())
	if errors.Is(err, http.ErrNoCookie) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return auth.Recaller(recaller), nil
}

func (sg *SessionGuard) RestoreAuth(ctx context.Context) error {
	ok, err := sg.restoreFromSession(ctx)
	if ok || err != nil {
		return err
	}

	_, err = sg.restoreFromRecaller(ctx)
	return err
}

func (sg *SessionGuard) restoreFromSession(ctx context.Context) (bool, error) {
	id, ok := sg.Session.Store.Get(sg.GetName())
	if ok {
		user, err := sg.Provider.RetrieveById(ctx, id)
		if errors.Is(err, auth.ErrUserNotFound) {
			return false, nil
		} else if err != nil {
			return false, err
		}
		sg.SetUser(user)
		return true, nil
	}
	return false, nil
}

func (sg *SessionGuard) restoreFromRecaller(ctx context.Context) (bool, error) {
	recaller, err := sg.GetRecaller()
	if err != nil {
		return false, err
	}

	if recaller != "" && recaller.Valid() {
		user, err := sg.Provider.RetrieveByToken(ctx, recaller.ID(), recaller.Token())
		if errors.Is(err, auth.ErrUserNotFound) {
			return false, nil
		} else if err != nil {
			return false, err
		}
		sg.SetUser(user)
		return true, nil
	}

	return false, nil
}

func (sg *SessionGuard) clearUserDataFromStorage() {
	sg.Session.Store.Remove(sg.GetName())
	sg.Cookie.Forget(sg.GetRecallerName())
}
