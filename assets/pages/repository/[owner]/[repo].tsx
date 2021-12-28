import Repository from "@/src/pages/repository";
import { useRouter } from "next/router";

export default function Page() {
  const router = useRouter();

  return (
    <Repository
      owner={router.query.owner as string}
      repo={router.query.repo as string}
    />
  );
}
