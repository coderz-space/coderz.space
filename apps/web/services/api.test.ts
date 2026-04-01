const mockRequest = jest.fn();
const mockCreate = jest.fn();
const mockPost = jest.fn();

jest.mock("axios", () => ({
  __esModule: true,
  default: {
    create: (...args: unknown[]) => mockCreate(...args),
    post: (...args: unknown[]) => mockPost(...args),
  },
}));

describe("api client", () => {
  beforeEach(() => {
    jest.resetModules();
    mockRequest.mockReset();
    mockCreate.mockReset();
    mockPost.mockReset();

    mockCreate.mockReturnValue({
      request: mockRequest,
      interceptors: {
        response: {
          use: jest.fn(),
        },
      },
    });
  });

  it("unwraps backend envelopes", async () => {
    mockRequest.mockResolvedValue({
      data: {
        success: true,
        data: { role: "mentee", accountStatus: "approved" },
      },
    });

    const { api } = await import("./api");

    await expect(api.get("/v1/app/context")).resolves.toEqual({
      role: "mentee",
      accountStatus: "approved",
    });
  });

  it("returns direct payloads without modification", async () => {
    mockRequest.mockResolvedValue({
      data: { status: "ok" },
    });

    const { api } = await import("./api");

    await expect(api.get("/health")).resolves.toEqual({ status: "ok" });
  });
});
