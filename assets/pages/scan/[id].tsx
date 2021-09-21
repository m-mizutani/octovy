import { useRouter } from "next/router";

function scan() {
  const router = useRouter();
  return <h1>Scan result: {router.query.id}</h1>;
}

export default scan;
