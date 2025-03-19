import useSWR from "swr";
import { Image } from "../types/image";

const fetcher = (url: string): Promise<Image[]> =>
  fetch(url).then((res) => res.json());

// 画像取得
export const useImages = () => {
  const { data, error } = useSWR("http://localhost:8080/get/images", fetcher, {
    suspense: true,
    revalidateOnFocus: false, // フォーカス時の再フェッチを無効
  });

  return {
    images: data || [],
    isLoading: !data && !error,
    isError: error,
  };
};
