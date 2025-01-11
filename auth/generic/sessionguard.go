package generic

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/wolftotem4/golava-core/auth"
	"github.com/wolftotem4/golava-core/auth/callback"
	"github.com/wolftotem4/golava-core/cookie"
	"github.com/wolftotem4/golava-core/hashing"
	"github.com/wolftotem4/golava-core/session"
	"github.com/wolftotem4/golava-core/util"
)

type SessionGuard struct {
	Name             string
	Session          *session.SessionManager
	Cookie           cookie.IEncryptableCookieManager
	Hasher           hashing.Hasher
	Request          *http.Request
	RememberDuration time.Duration
	Provider         auth.UserProvider
	RecallerIdMorph  auth.RecallerIdMorph

	// Listen to auth events.
	//
	// Recommend attaching this to an event emitter.
	//
	// Example:
	//
	// 	guard := &generic.SessionGuard{
	// 		Callbacks: callback.Listen(struct {
	// 			Attempting func(name string, credentials map[string]any, remember bool) error
	// 		}{
	// 			Attempting: func(name string, credentials map[string]any, remember bool) error {
	// 				// do something
	// 				return nil
	// 			},
	// 		}),
	// 	}
	//
	// See [github.com/wolftotem4/golava-core/auth/callback] for more information.
	Callbacks callback.Callbacks

	user            auth.Authenticatable
	viaRemember     bool
	currentRecaller *auth.Recaller
	userHash        string
}

func (sg *SessionGuard) User() auth.Authenticatable {
	return sg.user
}

func (sg *SessionGuard) UserHash() string {
	return sg.userHash
}

func (sg *SessionGuard) GetName() string {
	return fmt.Sprintf("login_%s", sg.Name)
}

func (sg *SessionGuard) GetRecallerName() string {
	return fmt.Sprintf("remember_%s", sg.Name)
}

func (sg *SessionGuard) SetUser(user auth.Authenticatable) error {
	return sg.setUser(context.TODO(), user, true, "")
}

func (sg *SessionGuard) setUser(ctx context.Context, user auth.Authenticatable, triggerAuthenticated bool, newhash string) error {
	sg.user = user

	if user == nil {
		sg.reset()
		return nil
	}

	if newhash != "" {
		sg.userHash = newhash
	} else {
		sg.userHash = user.GetAuthPassword()
	}

	if triggerAuthenticated && sg.Callbacks != nil {
		err := sg.Callbacks.Authenticated(ctx, sg.Name, user)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sg *SessionGuard) reset() {
	sg.user = nil
	sg.userHash = ""
	sg.viaRemember = false
	sg.setCurrentRecaller("")
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

func (sg *SessionGuard) validate(ctx context.Context, credentials map[string]any, shouldLogin ...auth.ShouldLogin) (auth.Authenticatable, bool, error) {
	user, err := sg.Provider.RetrieveByCredentials(ctx, credentials)
	if errors.Is(err, auth.ErrUserNotFound) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	valid, err := sg.Provider.ValidateCredentials(ctx, user, credentials)

	if sg.Callbacks != nil {
		err := sg.Callbacks.Validated(ctx, sg.Name, user)
		if err != nil {
			return user, false, err
		}
	}

	for _, validate := range shouldLogin {
		valid, err := validate(ctx, user)
		if err != nil {
			return user, false, err
		} else if !valid {
			return user, false, nil
		}
	}

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

func (sg *SessionGuard) Attempt(ctx context.Context, credentials map[string]any, remember bool, shouldLogin ...auth.ShouldLogin) (bool, error) {
	if sg.Callbacks != nil {
		err := sg.Callbacks.Attempting(ctx, sg.Name, credentials, remember)
		if err != nil {
			return false, err
		}
	}

	user, valid, err := sg.validate(ctx, credentials, shouldLogin...)
	if err != nil {
		return false, err
	} else if valid {
		newhash, err := sg.Provider.RehashPasswordIfRequired(ctx, user, credentials, false)
		if err != nil {
			return false, err
		}

		err = sg.login(ctx, user, remember, newhash)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	if sg.Callbacks != nil {
		err := sg.Callbacks.Failed(ctx, sg.Name, user)
		return false, err
	}

	return false, nil
}

// Log a user into the application without sessions or cookies.
func (sg *SessionGuard) Once(ctx context.Context, credentials map[string]any) (bool, error) {
	if sg.Callbacks != nil {
		err := sg.Callbacks.Attempting(ctx, sg.Name, credentials, false)
		if err != nil {
			return false, err
		}
	}

	user, valid, err := sg.validate(ctx, credentials)
	if err != nil {
		return false, err
	} else if !valid {
		return false, nil
	}

	newhash, err := sg.Provider.RehashPasswordIfRequired(ctx, user, credentials, false)
	if err != nil {
		return false, err
	}

	return true, sg.setUser(ctx, user, true, newhash)
}

// Log the given user ID into the application without sessions or cookies.
func (sg *SessionGuard) OnceUsingID(ctx context.Context, id any) (bool, error) {
	user, err := sg.Provider.RetrieveById(ctx, id)
	if err != nil {
		return false, err
	}

	return true, sg.setUser(ctx, user, true, "")
}

func (sg *SessionGuard) Login(ctx context.Context, user auth.Authenticatable, remember bool) error {
	return sg.login(ctx, user, remember, "")
}

func (sg *SessionGuard) login(ctx context.Context, user auth.Authenticatable, remember bool, newhash string) error {
	err := sg.updateSession(ctx, user.GetAuthIdentifier())
	if err != nil {
		return err
	}

	if remember {
		err = sg.ensureRememberTokenIsSet(ctx, user)
		if err != nil {
			return err
		}
		sg.createRecaller(user, newhash)
	}

	if sg.Callbacks != nil {
		err := sg.Callbacks.Login(ctx, sg.Name, user, remember)
		if err != nil {
			return err
		}
	}

	sg.Cookie.Encryption().Set(
		sg.Session.Name,
		sg.Session.Store.ID,
		cookie.WithMaxAge(int(sg.Session.Lifetime.Seconds())),
	)

	return sg.setUser(ctx, user, true, newhash)
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

func (sg *SessionGuard) createRecaller(user auth.Authenticatable, newhash string) {
	if newhash == "" {
		newhash = user.GetAuthPassword()
	}

	recaller := auth.NewRecallerString(
		user.GetAuthIdentifier(),
		user.GetRememberToken(),
		newhash,
	)
	sg.Cookie.Encryption().Set(
		sg.GetRecallerName(),
		recaller,
		cookie.WithMaxAge(int(sg.RememberDuration.Seconds())),
	)

	sg.setCurrentRecaller(auth.Recaller(recaller))
}

func (sg *SessionGuard) setCurrentRecaller(recaller auth.Recaller) {
	sg.currentRecaller = &recaller
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

	return sg.setUser(ctx, nil, false, "")
}

func (sg *SessionGuard) LogoutCurrentDevice(ctx context.Context) error {
	user := sg.User()

	sg.clearUserDataFromStorage()

	if sg.Callbacks != nil {
		err := sg.Callbacks.CurrentDeviceLogout(ctx, sg.Name, user)
		if err != nil {
			return err
		}
	}

	return sg.setUser(ctx, nil, false, "")
}

func (sg *SessionGuard) LogoutOtherDevices(ctx context.Context, password string) error {
	if sg.user == nil {
		return nil
	}

	newhash, err := sg.rehashUserPasswordForDeviceLogout(ctx, sg.user, password)
	if err != nil {
		return err
	}

	if newhash != "" {
		sg.userHash = newhash
	}

	if check, err := sg.checkRecaller(); err != nil {
		return err
	} else if check {
		sg.createRecaller(sg.user, newhash)
	}

	if sg.Callbacks != nil {
		return sg.Callbacks.OtherDeviceLogout(ctx, sg.Name, sg.user)
	}

	return nil
}

func (sg *SessionGuard) GetRecaller() (auth.Recaller, error) {
	if check, err := sg.checkRecaller(); err != nil {
		return "", err
	} else if check {
		return *sg.currentRecaller, nil
	}

	return "", nil
}

func (sg *SessionGuard) getRecaller() (auth.Recaller, error) {
	if sg.currentRecaller != nil {
		return *sg.currentRecaller, nil
	}

	value, err := sg.Cookie.Encryption().Get(sg.GetRecallerName())
	if err != nil {
		return "", nil
	}

	recaller := auth.Recaller(value)
	sg.setCurrentRecaller(recaller)
	return recaller, nil
}

func (sg *SessionGuard) RestoreAuth(ctx context.Context) error {
	ok, err := sg.restoreFromSession(ctx)
	if ok || err != nil {
		return err
	}

	_, err = sg.restoreFromRecaller(ctx)
	return err
}

func (sg *SessionGuard) ViaRemember() bool {
	return sg.viaRemember
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
		return true, sg.setUser(ctx, user, true, "")
	}
	return false, nil
}

func (sg *SessionGuard) restoreFromRecaller(ctx context.Context) (bool, error) {
	user, err := sg.getRecallerUser(ctx)
	if err != nil {
		return false, err
	}

	if user != nil {
		err := sg.setUser(ctx, user, true, "")
		if err != nil {
			return false, err
		}

		sg.viaRemember = true
		return true, nil
	}

	return false, nil
}

func (sg *SessionGuard) getRecallerUser(ctx context.Context) (auth.Authenticatable, error) {
	recaller, err := sg.getRecaller()
	if err != nil {
		return nil, err
	}

	if recaller != "" && recaller.Valid() {
		var id any
		if sg.RecallerIdMorph != nil {
			id, err = sg.RecallerIdMorph(recaller.ID())
			if err != nil {
				// ignore error
				return nil, nil
			}
		} else {
			id = recaller.ID()
		}

		user, err := sg.Provider.RetrieveByToken(ctx, id, recaller.Token())
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, nil
		}
		return user, err
	}

	return nil, nil
}

func (sg *SessionGuard) checkRecaller() (check bool, err error) {
	if sg.viaRemember {
		return true, nil
	} else if sg.user == nil {
		return false, nil
	}

	recaller, err := sg.getRecaller()
	if err != nil {
		return false, err
	}

	if recaller == "" || !recaller.Valid() {
		return false, nil
	}

	userId := sg.user.GetAuthIdentifier()
	if !recaller.MatchID(userId) {
		return false, nil
	}

	return true, nil
}

func (sg *SessionGuard) clearUserDataFromStorage() {
	sg.Session.Store.Remove(sg.GetName())
	sg.Cookie.Forget(sg.GetRecallerName())
}

func (sg *SessionGuard) rehashUserPasswordForDeviceLogout(ctx context.Context, user auth.Authenticatable, password string) (newhash string, err error) {
	check, err := sg.Hasher.Check(password, user.GetAuthPassword())
	if err != nil {
		return "", err
	}
	if !check {
		return "", auth.ErrPasswordMismatch
	}

	return sg.Provider.RehashPasswordIfRequired(ctx, user, map[string]any{user.GetAuthPasswordName(): password}, true)
}
