package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk"`
	Name      string    `bun:"type:varchar(100),notnull"`
	Email     string    `bun:"type:varchar(100),unique,notnull"`
	Password  string    `bun:"type:varchar(255),notnull" json:"-"`
	CreatedAt time.Time `bun:"type:timestamp,default:now(),notnull"`
	UpdatedAt time.Time `bun:"type:timestamp,default:now(),notnull"`
	DeletedAt time.Time `bun:"type:timestamp,soft_delete,nullzero" json:"-"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// HashPassword hashes the user's password with bcrypt
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password matches the hashed password
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
