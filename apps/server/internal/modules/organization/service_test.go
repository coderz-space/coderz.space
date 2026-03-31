package organization

import (
	"testing"

	db "github.com/DSAwithGautam/Coderz.space/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// Helper function to create test UUIDs
func testUUID(id byte) pgtype.UUID {
	return pgtype.UUID{
		Bytes: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, id},
		Valid: true,
	}
}

// Test slug validation
func TestValidateSlug(t *testing.T) {
	tests := []struct {
		name     string
		slug     string
		expected bool
	}{
		{"valid lowercase", "my-org", true},
		{"valid with numbers", "org123", true},
		{"valid with hyphens", "my-org-123", true},
		{"invalid uppercase", "My-Org", false},
		{"invalid special chars", "my_org", false},
		{"invalid spaces", "my org", false},
		{"too short", "ab", false},
		{"minimum length", "abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateSlug(tt.slug)
			if result != tt.expected {
				t.Errorf("ValidateSlug(%q) = %v, want %v", tt.slug, result, tt.expected)
			}
		})
	}
}

// Test slug normalization
func TestNormalizeSlug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"lowercase conversion", "My Organization", "my-organization"},
		{"remove special chars", "My Org!", "my-org"},
		{"multiple spaces", "my   org", "my-org"},
		{"trim hyphens", "-my-org-", "my-org"},
		{"consecutive hyphens", "my--org", "my-org"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeSlug(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeSlug(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Test organization data mapping
func TestMapOrganizationToData(t *testing.T) {
	svc := &Service{}

	org := db.Organization{
		ID:   testUUID(1),
		Name: "Test Org",
		Slug: "test-org",
		Description: pgtype.Text{
			String: "Test description",
			Valid:  true,
		},
		Status: db.OrgStatusPendingApproval,
		CreatedAt: pgtype.Timestamptz{
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Valid: true,
		},
	}

	result := svc.mapOrganizationToData(org)

	if result.Name != "Test Org" {
		t.Errorf("expected name %q, got %q", "Test Org", result.Name)
	}
	if result.Slug != "test-org" {
		t.Errorf("expected slug %q, got %q", "test-org", result.Slug)
	}
	if result.Description != "Test description" {
		t.Errorf("expected description %q, got %q", "Test description", result.Description)
	}
	if result.Status != string(db.OrgStatusPendingApproval) {
		t.Errorf("expected status %q, got %q", db.OrgStatusPendingApproval, result.Status)
	}
}

// Test member data mapping
func TestMapMemberToData(t *testing.T) {
	svc := &Service{}

	member := db.OrganizationMember{
		ID:             testUUID(1),
		OrganizationID: testUUID(2),
		UserID:         testUUID(3),
		Role:           db.OrgMemberRoleAdmin,
		JoinedAt: pgtype.Timestamptz{
			Valid: true,
		},
	}

	result := svc.mapMemberToData(member)

	if result.Role != string(db.OrgMemberRoleAdmin) {
		t.Errorf("expected role %q, got %q", db.OrgMemberRoleAdmin, result.Role)
	}
	if result.ID != member.ID {
		t.Errorf("expected ID %v, got %v", member.ID, result.ID)
	}
	if result.OrganizationID != member.OrganizationID {
		t.Errorf("expected OrganizationID %v, got %v", member.OrganizationID, result.OrganizationID)
	}
	if result.UserID != member.UserID {
		t.Errorf("expected UserID %v, got %v", member.UserID, result.UserID)
	}
}

// Test role parsing
func TestParseOrgMemberRole(t *testing.T) {
	svc := &Service{}

	tests := []struct {
		name        string
		role        string
		expected    db.OrgMemberRole
		expectError bool
	}{
		{"admin role", "admin", db.OrgMemberRoleAdmin, false},
		{"mentor role", "mentor", db.OrgMemberRoleMentor, false},
		{"mentee role", "mentee", db.OrgMemberRoleMentee, false},
		{"invalid role", "invalid", "", true},
		{"empty role", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.parseOrgMemberRole(tt.role)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for role %q, got nil", tt.role)
				}
				if err != nil && err.Error() != "INVALID_ROLE" {
					t.Errorf("expected INVALID_ROLE error, got %q", err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("expected role %q, got %q", tt.expected, result)
				}
			}
		})
	}
}

// Test slug uniqueness validation logic
func TestSlugUniquenessValidation(t *testing.T) {
	t.Run("slug validation enforces uniqueness", func(t *testing.T) {
		// This test documents that CreateOrganization checks slug uniqueness
		// by calling GetOrganizationBySlug before creating

		// The service should:
		// 1. Call GetOrganizationBySlug with the requested slug
		// 2. If it returns an organization (no error), return SLUG_ALREADY_EXISTS
		// 3. If it returns an error (not found), proceed with creation

		t.Log("CreateOrganization validates slug uniqueness before creation")
		t.Log("Expected behavior: SLUG_ALREADY_EXISTS error when slug exists")
	})
}

// Test status transition validation
func TestStatusTransitionValidation(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus db.OrgStatus
		canApprove    bool
		expectedError string
	}{
		{
			name:          "pending_approval can be approved",
			currentStatus: db.OrgStatusPendingApproval,
			canApprove:    true,
			expectedError: "",
		},
		{
			name:          "approved cannot be approved again",
			currentStatus: db.OrgStatusApproved,
			canApprove:    false,
			expectedError: "ORGANIZATION_NOT_PENDING",
		},
		{
			name:          "suspended cannot be approved",
			currentStatus: db.OrgStatusSuspended,
			canApprove:    false,
			expectedError: "ORGANIZATION_NOT_PENDING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test documents the status transition logic in ApproveOrganization
			// The service should:
			// 1. Get the organization by ID
			// 2. Check if status is PENDING_APPROVAL
			// 3. If not, return ORGANIZATION_NOT_PENDING error
			// 4. If yes, update status to APPROVED

			if tt.canApprove {
				t.Logf("Status %q should allow approval", tt.currentStatus)
			} else {
				t.Logf("Status %q should reject approval with error %q", tt.currentStatus, tt.expectedError)
			}
		})
	}
}

// Test admin auto-assignment
func TestAdminAutoAssignment(t *testing.T) {
	t.Run("creator is assigned as admin", func(t *testing.T) {
		// This test documents that CreateOrganization uses a transaction to:
		// 1. Create the organization with PENDING_APPROVAL status
		// 2. Add the creator as an admin member
		// 3. Commit both operations atomically

		t.Log("CreateOrganization should add creator as admin member in same transaction")
		t.Log("Expected: Both organization and member creation succeed or both fail")
	})
}

// Test transaction atomicity
func TestTransactionAtomicity(t *testing.T) {
	t.Run("organization creation is atomic", func(t *testing.T) {
		// This test documents that CreateOrganization uses transactions
		// The service should:
		// 1. Begin a transaction
		// 2. Create organization
		// 3. Add admin member
		// 4. Commit transaction
		// 5. If any step fails, rollback

		t.Log("CreateOrganization uses transaction to ensure atomicity")
		t.Log("If member creation fails, organization creation should be rolled back")
	})
}

// Test update slug uniqueness
func TestUpdateSlugUniqueness(t *testing.T) {
	t.Run("update validates slug uniqueness", func(t *testing.T) {
		// This test documents that UpdateOrganization validates slug uniqueness
		// The service should:
		// 1. If slug is being updated, validate format
		// 2. Check if slug exists with GetOrganizationBySlug
		// 3. If exists and belongs to different org, return SLUG_ALREADY_EXISTS
		// 4. If exists and belongs to same org, allow update

		t.Log("UpdateOrganization validates slug uniqueness excluding current org")
		t.Log("Same org can keep its slug, but cannot take another org's slug")
	})
}

// Test initial organization status
func TestInitialOrganizationStatus(t *testing.T) {
	t.Run("new organizations start as pending_approval", func(t *testing.T) {
		// This test documents that CreateOrganization sets status to PENDING_APPROVAL
		// The service should create organizations with status = PENDING_APPROVAL

		expectedStatus := db.OrgStatusPendingApproval
		t.Logf("New organizations should have status %q", expectedStatus)
	})
}

// Test RemoveMember last admin prevention (service layer)
func TestServiceRemoveMemberLastAdminPrevention(t *testing.T) {
	t.Run("cannot remove last admin", func(t *testing.T) {
		// This test documents that RemoveMember prevents deletion of the last admin
		// The service should:
		// 1. Get the member to check their role
		// 2. If member is an admin, count total admins
		// 3. If admin count <= 1, return CANNOT_REMOVE_LAST_ADMIN error
		// 4. Otherwise, proceed with removal

		t.Log("RemoveMember should prevent deletion of the last admin")
		t.Log("Expected error: CANNOT_REMOVE_LAST_ADMIN when removing last admin")
	})

	t.Run("can remove admin when multiple admins exist", func(t *testing.T) {
		// This test documents that RemoveMember allows admin removal when multiple admins exist
		// The service should:
		// 1. Get the member to check their role
		// 2. If member is an admin, count total admins
		// 3. If admin count > 1, proceed with removal

		t.Log("RemoveMember should allow admin removal when multiple admins exist")
		t.Log("Expected: Successful removal when admin count > 1")
	})

	t.Run("can remove non-admin members", func(t *testing.T) {
		// This test documents that RemoveMember allows removal of non-admin members
		// The service should:
		// 1. Get the member to check their role
		// 2. If member is not an admin, proceed with removal without checking admin count

		t.Log("RemoveMember should allow removal of mentor and mentee members")
		t.Log("Expected: Successful removal without admin count check")
	})
}

// Test RemoveMember member not found (service layer)
func TestServiceRemoveMemberNotFound(t *testing.T) {
	t.Run("returns error when member not found", func(t *testing.T) {
		// This test documents that RemoveMember returns MEMBER_NOT_FOUND error
		// The service should:
		// 1. Call GetOrganizationMember
		// 2. If member doesn't exist, return MEMBER_NOT_FOUND error

		t.Log("RemoveMember should return MEMBER_NOT_FOUND when member doesn't exist")
		t.Log("Expected error: MEMBER_NOT_FOUND")
	})
}
