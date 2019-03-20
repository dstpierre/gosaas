package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dstpierre/gosaas/model"
)

type Users struct {
	DB *sql.DB
}

func (u *Users) SignUp(email, password string) (*model.Account, error) {
	var accountID int64

	err := u.DB.QueryRow(`
		INSERT INTO gosaas_accounts(
			email, 
			stripe_id, 
			subscription_id, 
			plan, 
			is_yearly, 
			subscribed_on, 
			seats,
			is_active
		)
		VALUES($1, '', '', '', false, $2, 0, true)
		RETURNING id
	`, email, time.Now()).Scan(&accountID)
	if err != nil {
		return nil, err
	}

	_, err = u.DB.Exec(`
		INSERT INTO gosaas_users(account_id, email, password, token, role)
		VALUES($1, $2, $3, $4, $5)
	`, accountID, email, password, model.NewToken(accountID), model.RoleAdmin)
	if err != nil {
		return nil, err
	}

	return u.GetDetail(accountID)
}

func (u *Users) Auth(accountID int64, token string, pat bool) (*model.Account, *model.User, error) {
	token = fmt.Sprintf("%d|%s", accountID, token)

	user := &model.User{}
	row := u.DB.QueryRow("SELECT * FROM gosaas_users WHERE account_id = $1 AND token = $2", accountID, token)
	if err := u.scanUser(row, user); err != nil {
		return nil, nil, err
	}

	account, err := u.GetDetail(user.AccountID)
	if err != nil {
		return nil, nil, err
	}

	return account, user, nil
}

func (u *Users) GetDetail(id int64) (*model.Account, error) {
	account := &model.Account{}
	row := u.DB.QueryRow("SELECT * FROM gosaas_accounts WHERE id = $1", id)
	err := row.Scan(&account.ID,
		&account.Email,
		&account.StripeID,
		&account.SubscriptionID,
		&account.Plan,
		&account.IsYearly,
		&account.SubscribedOn,
		&account.Seats,
		&account.IsActive,
	)
	if err != nil {
		fmt.Println("error while scanning account")
		return nil, err
	}

	rows, err := u.DB.Query("SELECT * FROM gosaas_users WHERE account_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		if err := u.scanUser(rows, &user); err != nil {
			return nil, err
		}

		account.Users = append(account.Users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return account, nil
}

func (u *Users) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	row := u.DB.QueryRow("SELECT * FROM gosaas_users WHERE email = $1", email)

	if err := u.scanUser(row, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *Users) GetByStripe(stripeID string) (*model.Account, error) {
	var accountID int64
	row := u.DB.QueryRow("SELECT id FROM gosaas_accounts WHERE stripe_id = $1", stripeID)
	if err := row.Scan(&accountID); err != nil {
		return nil, err
	}

	return u.GetDetail(accountID)
}

func (u *Users) SetSeats(id int64, seats int) error {
	_, err := u.DB.Exec(`
		UPDATE gosaas_accounts SET
			seats = $2
		WHERE id = $1
	`, id, seats)
	return err
}

func (u *Users) ConvertToPaid(id int64, stripeID, subID, plan string, yearly bool, seats int) error {
	_, err := u.DB.Exec(`
		UPDATE gosaas_accounts SET
			stripe_id = $2,
			subscription_id = $3,
			subscribed_on = $4,
			plan = $5,
			seats = $6,
			is_yearly = $7
		WHERE id = $1
	`, id, stripeID, subID, time.Now(), plan, seats, yearly)
	return err
}

func (u *Users) ChangePlan(id int64, plan string, yearly bool) error {
	_, err := u.DB.Exec(`
		UPDATE gosaas_accounts SET
			plan = $2,
			is_yearly = $3
		WHERE id = $1
	`, id, plan, yearly)
	return err
}

func (u *Users) Cancel(id int64) error {
	_, err := u.DB.Exec(`
		UPDATE gosaas_accounts SET
			subscription_id = '',
			plan = '',
			is_yearly = false
		WHERE id = $1
	`, id)
	return err
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func (u *Users) AddToken(accountID, userID int64, name string) (*model.AccessToken, error) {
	return nil, fmt.Errorf("not implemented")
}

func (u *Users) RemoveToken(accountID, userID, tokenID int64) error {
	return fmt.Errorf("not implemented")
}

func (u *Users) scanUser(rows scanner, user *model.User) error {

	return rows.Scan(&user.ID,
		&user.AccountID,
		&user.Email,
		&user.Password,
		&user.Token,
		&user.Role,
	)
}
