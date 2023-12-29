import useSWR from "swr";
import { PaginationResult } from "../types/pagination";
import { Submission } from "../types/submission";
import { fetcherWithAuth } from "./fetcher";

export function usePaginationSubmission(
  page: number,
  limit: number,
  sort: string,
  order: string,
  challengeID: number,
  userID: number
) {
  let url = `/submission/pagination?page=${page}&limit=${limit}&sort=${sort}&order=${order}`;
  if (challengeID) {
    url += `&challengeId=${challengeID}`;
  }
  if (userID) {
    url += `&userId=${userID}`;
  }

  const { data, error, isLoading } = useSWR(() => url, fetcherWithAuth);

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
