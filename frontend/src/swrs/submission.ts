import useSWR from "swr";
import { PaginationResult } from "../types/pagination";
import { Submission } from "../types/submission";
import { fetcherWithAuth } from "./fetcher";

export function usePaginationSubmission(
  page: number,
  limit: number,
  sort: string,
  order: string,
  search: string,
  challengeID?: any,
  userID?: any
) {
  const { data, error, isLoading } = useSWR(() => {
    const params = new URLSearchParams({
      page: String(page),
      limit: String(limit),
      sort,
      order,
      search,
    });
    if (challengeID) {
      params.append("challenge_id", String(challengeID));
    }
    if (userID) {
      params.append("user_id", String(userID));
    }

    return `/submission/pagination?${params.toString()}`;
  }, fetcherWithAuth);

  return {
    data: data as PaginationResult<Submission>,
    isLoading,
    isError: error,
  };
}

export function useSubmission(submissionID: any) {
  const { data, error, isLoading } = useSWR(
    () => `/submission/get/${submissionID}`,
    fetcherWithAuth
  );

  return {
    data: data as Submission,
    isLoading,
    isError: error,
  };
}
