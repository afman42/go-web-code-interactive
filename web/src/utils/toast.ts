//https://github.com/afman42/simple-image-uploader-sveltekit/blob/main/src/lib/store/state.ts
import { writable, derived } from "svelte/store";
export type IToast = {
  message: string;
  type: TToastStatus;
  id: number;
  timeout: number;
};
export type TToastStatus = "success" | "warning" | "error" | "info";

export const state = writable<boolean>(true);

function storeToast() {
  const _toasts = writable<IToast[]>([]);

  function send(
    message: string,
    type: TToastStatus = "success",
    timeout: number,
  ) {
    _toasts.update(
      //@ts-ignore
      (state: IToast[]) =>
        (state = [
          ...state,
          { id: Math.random() * 1000, type, message, timeout },
        ]),
    );
  }

  const toasts = derived(_toasts, ($_toasts, set) => {
    set($_toasts);
    if ($_toasts.length > 0) {
      const timer = setTimeout(() => {
        _toasts.update((state) => {
          state.shift();
          return state;
        });
      }, $_toasts[0].timeout);
      return () => {
        clearTimeout(timer);
      };
    }
  });
  const { subscribe } = toasts;

  return {
    stateToast: toasts,
    subscribe,
    send,
    error: (msg: string, timeout: number) => send(msg, "error", timeout),
    warning: (msg: string, timeout: number) => send(msg, "warning", timeout),
    success: (msg: string, timeout: number) => send(msg, "success", timeout),
    info: (msg: string, timeout: number) => send(msg, "info", timeout),
  };
}
export const toasts = storeToast();
