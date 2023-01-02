package db

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gitlab.boomerangapp.ir/back/pg/schema/entities"
	"gitlab.boomerangapp.ir/back/pg/types"
	"gitlab.boomerangapp.ir/back/pg/types/consts"
	utilsConsts "gitlab.boomerangapp.ir/back/utils/consts"
)

func LoadWalletBalance(ctx context.Context, userID uint) (int64, error) {

	wallet, err := entities.Wallets(entities.WalletWhere.UserID.EQ(userID),
		qm.Select(entities.WalletColumns.Balance)).One(ctx, DB)
	if err != nil {
		return 0, err
	}

	return wallet.Balance, nil
}

func LoadDailyDepositSum(ctx context.Context, userID uint) (uint64, error) {

	row := DB.QueryRowContext(ctx, "SELECT SUM(amount) FROM payments WHERE status = 1 "+
		"AND (confirm_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL 1 DAY) "+
		"AND user_id = ?", userID)
	if row.Err() != nil {
		return 0, row.Err()
	}

	var sumAmount *uint64
	err := row.Scan(&sumAmount)
	if err != nil {
		return 0, err
	}

	if sumAmount == nil {
		return 0, nil
	}
	return *sumAmount, nil
}

func LoadMaxRules(ctx context.Context, userID uint) (uint, uint, error) {

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}

	user, err := entities.FindUser(ctx, tx, userID, entities.UserColumns.UserLevelID)
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	level_accesses, err := entities.UserLevelAccesses(
		entities.UserLevelAccessWhere.LevelID.EQ(user.UserLevelID),
		entities.UserLevelAccessWhere.AccessID.IN([]uint8{1, 2})).All(ctx, tx)
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	var maxBalance, maxDeposit uint
	for _, la := range level_accesses {
		lav, _ := strconv.ParseUint(la.AccessValue, 10, 32)
		if la.AccessID == 1 {
			maxBalance = uint(lav)
			continue
		}
		maxDeposit = uint(lav)
	}

	err = tx.Commit()
	if err != nil {
		return 0, 0, err
	}

	return maxBalance, maxDeposit, nil
}

func IncreaseWalletBalance(ctx context.Context, tx *sql.Tx, walletID, amount uint64) error {

	_, err := tx.ExecContext(ctx, "UPDATE wallets SET balance = balance + ? WHERE wallet_id = ?",
		amount, walletID)
	if err != nil {
		return err
	}

	return nil
}

func LoadMinDeposit(ctx context.Context) (uint, error) {

	var result uint
	err := DB.QueryRowContext(ctx, "SELECT int_value FROM settings WHERE s_key = 'MIN_DEPOSIT'").
		Scan(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func LoadPaymentsByUserID(ctx context.Context, userID uint, isProducer bool, limit, offset int) (entities.PaymentSlice, int64, error) {

	where := entities.PaymentWhere.UserID.EQ(null.Uint{Uint: userID, Valid: true})
	if isProducer {
		where = entities.PaymentWhere.ProducerID.EQ(null.Uint{Uint: userID, Valid: true})
	}

	payments, err := entities.Payments(where, qm.Limit(limit), qm.Offset(offset)).All(ctx, DB)
	if err != nil {
		return nil, 0, err
	}

	total, err := entities.Payments(where).Count(ctx, DB)
	if err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

func UpdateStatus(ctx context.Context, tx *sql.Tx, id uint, status uint8) error {

	_, err := entities.Payments(entities.PaymentWhere.ID.EQ(id)).UpdateAll(ctx, tx, entities.M{
		entities.PaymentColumns.Status: status,
	})
	if err != nil {
		return err
	}

	return nil
}

func UpdateVerify(ctx context.Context, tx *sql.Tx, id, amount uint, cardMask, resp string) error {

	_, err := entities.Payments(entities.PaymentWhere.ID.EQ(id)).UpdateAll(ctx, tx, entities.M{
		entities.PaymentColumns.Amount:     amount,
		entities.PaymentColumns.VerifyResp: resp,
		entities.PaymentColumns.CardMask:   cardMask,
	})
	if err != nil {
		return err
	}

	return nil
}

func UpdateConfirm(ctx context.Context, tx *sql.Tx, id uint, resp string) error {

	_, err := entities.Payments(entities.PaymentWhere.ID.EQ(id)).UpdateAll(ctx, tx, entities.M{
		entities.PaymentColumns.ConfirmDate: time.Now(),
		entities.PaymentColumns.ConfirmResp: resp,
	})
	if err != nil {
		return err
	}

	return nil
}

func LoadPayment(ctx context.Context, tx *sql.Tx, id uint) (*entities.Payment, error) {

	var exec boil.ContextExecutor = tx
	if tx == nil {
		exec = DB
	}

	payment, err := entities.FindPayment(ctx, exec, id)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func UpdateToken(ctx context.Context, id uint, token string) error {

	_, err := entities.Payments(entities.PaymentWhere.ID.EQ(id)).
		UpdateAll(ctx, DB, entities.M{entities.PaymentColumns.PSPToken: token})
	if err != nil {
		return err
	}

	return nil
}

func CreatePayment(ctx context.Context, pr types.PaymentRequest, role, ip string, id uint) (uint, error) {

	payment := &entities.Payment{
		WalletID:     pr.WalletID,
		Amount:       pr.Amount,
		ActionType:   consts.ActionType,
		ActionID:     strconv.FormatUint(pr.WalletID, 10),
		Callback:     pr.Callback,
		GatewayID:    pr.GatewayID,
		AutoRedirect: pr.AutoRedirect,
		PSP:          pr.PSP,
		IP:           ip,
	}
	if role == utilsConsts.RoleProducer {
		payment.ProducerID = null.Uint{
			Uint:  id,
			Valid: true,
		}
	} else {
		payment.UserID = null.Uint{
			Uint:  id,
			Valid: true,
		}
	}
	err := payment.Insert(ctx, DB, boil.Infer())
	if err != nil {
		return 0, err
	}

	return payment.ID, nil
}
