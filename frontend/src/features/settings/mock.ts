// Mock sessions matching GET /account/sessions until the settings API is wired.
export type MockSession = {
  id: string;
  device: string;
  ip: string;
  current: boolean;
  createdAt: string;
  expiresAt: string;
};

export const MOCK_SESSIONS: MockSession[] = [
  {
    id: "sess_current",
    device: "Chrome · macOS",
    ip: "49.36.118.4",
    current: true,
    createdAt: "Jul 27, 2026",
    expiresAt: "Aug 03, 2026",
  },
  {
    id: "sess_2",
    device: "Safari · iPhone",
    ip: "39.45.12.88",
    current: false,
    createdAt: "Jul 25, 2026",
    expiresAt: "Aug 01, 2026",
  },
  {
    id: "sess_3",
    device: "Firefox · Linux",
    ip: "103.22.9.14",
    current: false,
    createdAt: "Jul 20, 2026",
    expiresAt: "Jul 27, 2026",
  },
];
