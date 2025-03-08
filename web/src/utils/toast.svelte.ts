type IToast = {
  message: string;
  type: TToastStatus;
  id: number;
  timeout: number;
};
type TToastStatus = "success" | "warning" | "error" | "info";

let _toasts = $state<IToast[]>([]);

function send(
  message: string,
  type: TToastStatus = "success",
  timeout: number,
) {
  _toasts.push({
    ..._toasts,
    id: Math.random() * 1000,
    type,
    message,
    timeout,
  });
}

const error = (msg: string, timeout: number) => send(msg, "error", timeout);
const warning = (msg: string, timeout: number) => send(msg, "warning", timeout);
const success = (msg: string, timeout: number) => send(msg, "success", timeout);
const info = (msg: string, timeout: number) => send(msg, "info", timeout);

export const useToast = () => {
  $effect(() => {
    if (_toasts.length > 0) {
      const timer = setTimeout(() => {
        //Remove first
        _toasts.shift();
      }, _toasts[0].timeout);
      return () => {
        clearTimeout(timer);
      };
    }
  });
  return {
    state: _toasts,
    error,
    warning,
    info,
    success,
  };
};
