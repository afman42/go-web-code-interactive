const envAPi = import.meta.env.VITE_API_URL;
export async function fetchApiPost<T>(
  payload: any = "",
  path: string,
): Promise<T> {
  const response = await fetch(envAPi + path, {
    method: "POST",
    body: JSON.stringify(payload),
    headers: {
      "Content-Type": "application/json;charset=utf-8",
    },
  });
  if (response.ok) return await response.json();
  const errorMesssage: string = await response.text();
  return Promise.reject(errorMesssage);
}
