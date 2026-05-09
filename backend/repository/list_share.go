package repository

import (
	"crypto/rand"
	"math/big"
	"time"
	"todo-list/backend/database"
	"todo-list/backend/models"
)

func generateCode() string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	code := make([]byte, 8)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		code[i] = chars[n.Int64()]
	}
	return string(code)
}

func CreateShare(listID uint, permission string) (*models.ListShare, error) {
	// If share already exists, return it
	var existing models.ListShare
	if database.DB.Where("list_id = ?", listID).First(&existing).Error == nil {
		return &existing, nil
	}
	s := &models.ListShare{ListID: listID, Code: generateCode(), Permission: permission}
	if err := database.DB.Create(s).Error; err != nil {
		// Retry once on code collision
		s.Code = generateCode()
		if err2 := database.DB.Create(s).Error; err2 != nil {
			return nil, err2
		}
	}
	return s, nil
}

func FindShareByListID(listID uint) (*models.ListShare, error) {
	var s models.ListShare
	err := database.DB.Where("list_id = ?", listID).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func FindShareByCode(code string) (*models.ListShare, error) {
	var s models.ListShare
	err := database.DB.Where("code = ?", code).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func DeleteShare(listID uint) error {
	s, err := FindShareByListID(listID)
	if err != nil {
		return err
	}
	database.DB.Where("share_id = ?", s.ID).Delete(&models.ListShareMember{})
	return database.DB.Delete(s).Error
}

func AddShareMember(shareID uint, userID uint) error {
	// Idempotent: if already a member, do nothing
	var count int64
	database.DB.Model(&models.ListShareMember{}).Where("share_id = ? AND user_id = ?", shareID, userID).Count(&count)
	if count > 0 {
		return nil
	}
	m := &models.ListShareMember{ShareID: shareID, UserID: userID, JoinedAt: time.Now()}
	return database.DB.Create(m).Error
}

func FindShareMembers(shareID uint) ([]models.ListShareMember, error) {
	var members []models.ListShareMember
	err := database.DB.Where("share_id = ?", shareID).Find(&members).Error
	return members, err
}

// FindSharedListIDs returns list IDs where the user is a share member
func FindSharedListIDs(userID uint) []uint {
	var ids []uint
	database.DB.Model(&models.ListShareMember{}).
		Joins("JOIN list_shares ON list_shares.id = list_share_members.share_id").
		Where("list_share_members.user_id = ?", userID).
		Pluck("list_shares.list_id", &ids)
	return ids
}

// GetUserPermission returns "" if no access, "view" or "edit" if member
func GetUserPermission(listID uint, userID uint) string {
	var s models.ListShare
	if database.DB.Where("list_id = ?", listID).First(&s).Error != nil {
		return ""
	}
	var count int64
	database.DB.Model(&models.ListShareMember{}).Where("share_id = ? AND user_id = ?", s.ID, userID).Count(&count)
	if count == 0 {
		return ""
	}
	return s.Permission
}
