package errorHandler

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
)

const (
	UniqueConstraint      = "23505"
	ForeignKeyConstraint  = "23503"
	NotNullConstraint     = "23502"
	CheckConstraint       = "23514"
	DataTypeMismatches    = "22P02"
	InvalidAmountOfParams = "42601"
	LostConnection        = "08003"
	ShutdownPGServer      = "57P01"
)

func OnHandleError(err error, op string) error {
	if pgErr := err.(*pq.Error); errors.As(err, &pgErr) {
		if pgErr.Code == UniqueConstraint {
			return fmt.Errorf("Alias already exists  %s: %w", op, err)
		} else if pgErr.Code == ForeignKeyConstraint {
			return fmt.Errorf("Foreign key constraint  %s: %w", op, err)
		} else if pgErr.Code == NotNullConstraint {
			return fmt.Errorf("Not null constraint  %s: %w", op, err)
		} else if pgErr.Code == CheckConstraint {
			return fmt.Errorf("Check constraint  %s: %w", op, err)
		} else if pgErr.Code == DataTypeMismatches {
			return fmt.Errorf("Data type mismatches  %s: %w", op, err)
		} else if pgErr.Code == InvalidAmountOfParams {
			return fmt.Errorf("Invalid amount of parameters  %s: %w", op, err)
		} else if pgErr.Code == LostConnection {
			return fmt.Errorf("Lost connection to server  %s: %w", op, err)
		} else if pgErr.Code == ShutdownPGServer {
			return fmt.Errorf("Shutdown postgres server  %s: %w", op, err)
		}
		fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
