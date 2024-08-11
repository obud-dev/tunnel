export const request = async <T>(
  url: string,
  options?: RequestInit
): Promise<T> => {
  const response = await fetch(url, options);
  if (!response.ok) {
    throw new Error(response.statusText);
  }

  const { code, data, msg } = await response.json();
  if (code === 0) {
    return data;
  }
  return Promise.reject(msg);
};

export interface Tunnel {
  id: string;
  name: string;
  status: string;
  uptime: number;
}

export interface Route {
  id: string;
  tunnel_id: string;
  hostname: string;
  prefix: string;
  target: string;
  protocol: string;
}
