import useSWR from "swr";
import { fetcherWithAuth } from "./fetcher";

export function useSubmission(
  page: number,
  limit: number,
  sort: string,
  order: string,
  challengeId: string,
  userId: string
) {
  let url = `/challenge/pagination?page=${page}&limit=${limit}&sort=${sort}&order=${order}`;
  if (challengeId) {
    url += `&challengeId=${challengeId}`;
  }
  if (userId) {
    url += `&userId=${userId}`;
  }

  const { data, error, isLoading } = useSWR(url, fetcherWithAuth);

  return {
    user: data,
    isLoading,
    isError: error,
  };
}
