import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import MenteeSignUpCard from "./MenteeSignUpCard";
import { registerMentee } from "@/services/auth";

jest.mock("@/services/auth", () => ({
  registerMentee: jest.fn(),
}));

describe("MenteeSignUpCard", () => {
  beforeEach(() => {
    (registerMentee as jest.Mock).mockReset();
  });

  it("shows the pending approval message after signup", async () => {
    (registerMentee as jest.Mock).mockResolvedValue({
      id: "request-1",
      firstName: "Alice",
      lastName: "Example",
      username: "alice",
      email: "alice@example.com",
      signedUpAt: new Date().toISOString(),
      status: "pending",
    });

    const user = userEvent.setup();
    render(<MenteeSignUpCard role="mentee" onClose={jest.fn()} onBackToLogin={jest.fn()} />);

    await user.type(screen.getByPlaceholderText("First Name"), "Alice");
    await user.type(screen.getByPlaceholderText("Last Name"), "Example");
    await user.type(screen.getByPlaceholderText("Username"), "alice");
    await user.type(screen.getByPlaceholderText("Email"), "alice@example.com");
    await user.type(screen.getByPlaceholderText("Set Password (min 8 chars)"), "password123");
    await user.type(screen.getByPlaceholderText("Confirm Password"), "password123");
    await user.click(screen.getByRole("button", { name: "Sign Up" }));

    await waitFor(() => {
      expect(registerMentee).toHaveBeenCalledWith({
        firstName: "Alice",
        lastName: "Example",
        username: "alice",
        email: "alice@example.com",
        password: "password123",
      });
    });

    expect(
      screen.getByText("Sign-up request submitted. A mentor will approve your account before you can log in.")
    ).toBeInTheDocument();
  });
});
