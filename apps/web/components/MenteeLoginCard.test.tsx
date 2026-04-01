import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import MenteeLoginCard from "./MenteeLoginCard";
import { loginMenteeByEmail } from "@/services/auth";

const push = jest.fn();

jest.mock("next/navigation", () => ({
  useRouter: () => ({
    push,
  }),
}));

jest.mock("@/services/auth", () => ({
  loginMenteeByEmail: jest.fn(),
}));

describe("MenteeLoginCard", () => {
  beforeEach(() => {
    push.mockReset();
    (loginMenteeByEmail as jest.Mock).mockReset();
  });

  it("routes approved mentees to their dashboard", async () => {
    (loginMenteeByEmail as jest.Mock).mockResolvedValue({
      auth: {
        accessToken: "access",
        refreshToken: "refresh",
        user: {
          id: "1",
          name: "Alice",
          email: "alice@example.com",
          emailVerified: true,
        },
      },
      context: {
        role: "mentee",
        accountStatus: "approved",
        user: {
          id: "1",
          name: "Alice Example",
          firstName: "Alice",
          lastName: "Example",
          username: "alice",
          email: "alice@example.com",
        },
      },
    });

    const user = userEvent.setup();
    render(<MenteeLoginCard role="mentee" onClose={jest.fn()} onSignUp={jest.fn()} />);

    await user.type(screen.getByPlaceholderText("Email address"), "alice@example.com");
    await user.type(screen.getByPlaceholderText("Password"), "password123");
    await user.click(screen.getByRole("button", { name: "Log In" }));

    await waitFor(() => {
      expect(loginMenteeByEmail).toHaveBeenCalledWith("alice@example.com", "password123");
      expect(push).toHaveBeenCalledWith("/mentee-dashboard/alice");
    });
  });
});
