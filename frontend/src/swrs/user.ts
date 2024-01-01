import useSWR from "swr";
import { PaginationResult } from "../types/pagination";
import { User } from "../types/user";
import { fetcherWithAuth } from "./fetcher";

export function usePaginationUser(
  page: number,
  limit: number,
  order: string,
  sort: string,
  search: string
) {
  const { data, error, isLoading } = useSWR(() => {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
      sort,
      order,
      search,
    });
    return `/user/pagination?${params.toString()}`;
  }, fetcherWithAuth);

  return {
    data: data as PaginationResult<User>,
    isLoading,
    isError: error,
  };
}
