import Head from "next/head";

export default function CommonHead({ title = "Octovy" }) {
  const description = "GitHub repository scanner";
  const bodyStyle = `body{margin: 0}`;
  return (
    <Head>
      <title>{title}</title>
      <meta property="description" content={description} />
      <meta property="og:title" content={title} />
      <meta property="og:description" content={description} />
      <link
        href="https://fonts.googleapis.com/css2?family=Kanit&display=swap"
        rel="stylesheet"
      />
      <style type="text/css">{bodyStyle}</style>
    </Head>
  );
}
