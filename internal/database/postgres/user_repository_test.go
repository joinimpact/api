package postgres_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/joinimpact/api/internal/database/postgres"
	"github.com/joinimpact/api/internal/models"
)

var sampleUsers []models.User = []models.User{
	{
		Model: models.Model{
			ID: 12738165059,
		},
		Active:        true,
		FirstName:     "Yury",
		LastName:      "Orlovskiy",
		Password:      "$2b$10$//DXiVVE59p7G5k/4Klx/ezF7BI42QZKmoOD0NDvUuqxRE5bFFBLy",
		Email:         "yury@joinimpact.org",
		EmailVerified: true,
	},
}

type gormTest struct {
	db   *gorm.DB
	mock sqlmock.Sqlmock
}

func initializeGorm() (*gormTest, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}

	gorm, err := gorm.Open("postgres", db)
	if err != nil {
		return nil, err
	}

	// Enable log mode to see database queries in the console.
	gorm.LogMode(true)

	return &gormTest{
		gorm,
		mock,
	}, nil
}

// TestCreate tests the Create function of the repository.
func TestCreate(t *testing.T) {
	// Initialize the mock database and ORM.
	test, err := initializeGorm()
	if err != nil {
		t.Fatal("error initializing mock database:", err)
	}
	defer test.db.Close()

	// Get the table name for the User model.
	usersTableName := test.db.NewScope(&models.User{}).TableName()

	// Create a repository around our mock database.
	repo := postgres.NewUserRepository(test.db, nil)

	// Attempt to add the test users into the database.
	for _, user := range sampleUsers {
		// Expect the beginning of a transaction.
		test.mock.ExpectBegin()

		// Expect a query inserting a user into the database.
		test.mock.ExpectQuery(fmt.Sprintf("INSERT INTO \"%s\"", usersTableName)).
			WithArgs(
				user.ID,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				user.Active,
				user.Email,
				user.EmailVerified,
				user.Password,
				user.ProfilePicture,
				user.FirstName,
				user.LastName,
				sqlmock.AnyArg(),
				user.ZIPCode,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(user.ID))

		// Expect a commit to the database.
		test.mock.ExpectCommit()

		// Run the Create function to create the user.
		err := repo.Create(user)
		if err != nil {
			t.Error("error creating user:", err)
		}
	}

	// Make sure that all expectations were met.
	if err := test.mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestFindByID tests the FindByID function of the repository.
func TestFindByID(t *testing.T) {
	// Initialize the mock database and ORM.
	test, err := initializeGorm()
	if err != nil {
		t.Fatal("error initializing mock database:", err)
	}
	defer test.db.Close()

	// Create a repository around our mock database.
	repo := postgres.NewUserRepository(test.db, nil)

	// Expect the select query from the repository.
	test.mock.ExpectQuery("SELECT").
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(sampleUsers[0].ID),
		)

	// Find the user by ID 12738165059.
	_, err = repo.FindByID(12738165059)
	if err != nil {
		t.Error("could not find user:", err)
	}

	// Make sure that all expectations were met.
	if err := test.mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
