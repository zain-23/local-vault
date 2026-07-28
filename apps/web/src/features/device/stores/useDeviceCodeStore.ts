import { create } from "zustand";

interface DeviceCodeState {
  userCode: string | null;
  setUserCode: (code: string) => void;
  clear: () => void;
}

// Holds the user code between the submit step and the approval screen, so it
// never rides in the URL (URLs land in history, logs, and referrer headers).
// Intentionally in-memory — NOT persisted: a refresh wipes it, so the approval
// screen falls back to "no code" and the user re-enters it from their terminal.
// It survives client-side navigation, so the login redirect round-trip is fine.
export const useDeviceCodeStore = create<DeviceCodeState>((set) => ({
  userCode: null,
  setUserCode: (code) => set({ userCode: code }),
  clear: () => set({ userCode: null }),
}));
