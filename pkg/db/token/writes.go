package token

import (
	"context"
	"fmt"
	"homeserver/pkg/interfaces"
)

func CreateToken(ctx context.Context, db interfaces.DatabaseInterface, userID int, token string) error {
	_, err := db.Exec(ctx, `UPDATE users SET auth_token = $1 WHERE id = $2`, token, userID)
	if err != nil {
		return fmt.Errorf("error creating token: %v", err)
	}
	return nil
}

func DeleteToken(ctx context.Context, db interfaces.DatabaseInterface, token string) error {
	_, err := db.Exec(ctx, `UPDATE users SET auth_token = NULL WHERE auth_token = $1`, token)
	if err != nil {
		return fmt.Errorf("error deleting token: %v", err)
	}
	return nil
}
