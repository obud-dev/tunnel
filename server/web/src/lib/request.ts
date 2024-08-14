export const request = async <T> (
  url: string,
  options?: RequestInit
): Promise<Resp<T>> => {
  const response = await fetch(url, options);
  if (!response.ok) {
    Promise.reject(response.statusText);
  }

  const resp = await response.json();
  return Promise.resolve(resp);
};

export interface Resp<T> {
  code: number;
  data: T;
  msg: string;
}

export interface Tunnel {
  id: string;
  name: string;
  token: string;
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
