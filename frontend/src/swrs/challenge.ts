import useSWR from "swr";
import { Challenge } from "../types/challenge";
import { PaginationResult } from "../types/pagination";
import { fetcherWithAuth } from "./fetcher";

export function usePaginationChallenge(
  page: number,
  limit: number,
  order: string,
  sort: string,
  search: string
) {
  const { data, error, isLoading } = useSWR(() => {
    const params = new URLSearchParams({
      page: String(page),
      limit: String(limit),
      sort,
      order,
      search,
    });

    return `/challenge/pagination?${params.toString()}`;
  }, fetcherWithAuth);

  return {
    data: data as PaginationResult<Challenge>,
    isLoading,
    isError: error,
  };
}

export function useChallenge(challengeID: any) {
  const { data, error, isLoading } = useSWR(
    () => `/challenge/get/${challengeID}`,
    fetcherWithAuth
  );

  return {
    data: data as Challenge,
    isLoading,
    isError: error,
  };
}
