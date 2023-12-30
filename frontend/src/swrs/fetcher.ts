const API_URL = import.meta.env.VITE_API_URL;

export const fetcher = (url: string) => {
  return fetch(API_URL + url).then((res) => res.json());
};

export const fetcherWithAuth = (url: string) => {
  const token = localStorage.getItem("accessToken");

  return fetch(API_URL + url, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  }).then((res) => {
    if (res.status === 200) {
      return res.json();
    }
    throw new Error("Error fetching data " + res.status);
  });
};
