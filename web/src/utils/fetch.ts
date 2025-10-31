const envAPI = import.meta.env.VITE_API_URL || "/";

export async function fetchApiPost<T>(
  payload: any = "",
  path: string,
): Promise<T> {
  // Validate payload before sending
  if (!validatePayload(payload)) {
    throw new Error(
      "Invalid payload: missing required fields or invalid types",
    );
  }

  const response = await fetch(envAPI + path, {
    method: "POST",
    body: JSON.stringify(payload),
    headers: {
      "Content-Type": "application/json;charset=utf-8",
      "X-Requested-With": "XMLHttpRequest",
      Accept: "application/json",
    },
  });

  // Validate response content type to prevent XSS
  const contentType = response.headers.get("content-type");
  if (!contentType || !contentType.includes("application/json")) {
    throw new Error("Invalid response type received from server");
  }

  if (response.ok) {
    const data = await response.json();
    return sanitizeResponse(data) as T;
  } else {
    // Handle different error status codes
    if (response.status >= 400 && response.status < 500) {
      const errorText = await response.text();
      return Promise.reject(errorText);
    } else {
      throw new Error(`Server error: ${response.status}`);
    }
  }
}

function validatePayload(payload: any): boolean {
  // Validate payload structure and content
  if (typeof payload !== "object" || payload === null) return false;
  if (typeof payload.txt !== "string") return false;
  if (typeof payload.lang !== "string") return false;
  if (typeof payload.type !== "string") return false;

  // Additional validation
  if (payload.txt.length > 10000) return false; // 10KB limit
  if (!["node", "php", "go"].includes(payload.lang)) return false;
  if (!["repl", "stq"].includes(payload.type)) return false;

  return true;
}

function sanitizeResponse(data: any): any {
  // Sanitize any output that will be displayed to prevent XSS
  if (data.out && typeof data.out === "string") {
    data.out = sanitizeHtml(data.out);
  }
  if (data.errout && typeof data.errout === "string") {
    data.errout = sanitizeHtml(data.errout);
  }

  return data;
}

function sanitizeHtml(html: string): string {
  // Basic HTML sanitization to prevent XSS
  return html
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#x27;");
}
