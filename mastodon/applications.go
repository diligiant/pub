package mastodon

import "github.com/jmoiron/sqlx"

type applications struct {
	db *sqlx.DB
}

func (a *applications) create(app *Application) error {
	result, err := a.db.NamedExec(`INSERT INTO applications (name, website, redirect_uri, client_id, client_secret, vapid_key) VALUES (:name, :website, :redirect_uri, :client_id, :client_secret, :vapid_key)`, app)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	app.ID = int(id)
	return nil
}

func (a *applications) findByClientID(clientID string) (*Application, error) {
	app := &Application{}
	err := a.db.QueryRowx(`SELECT * FROM applications WHERE client_id = ?`, clientID).StructScan(app)
	return app, err
}