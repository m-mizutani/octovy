import Document, { Html, Head, Main, NextScript } from "next/document";

class MyDocument extends Document {
  static async getInitialProps(ctx) {
    const initialProps = await Document.getInitialProps(ctx);
    return { ...initialProps };
  }

  render() {
    return (
      <Html>
        <Head>
          <link
            href="https://fonts.googleapis.com/css2?family=Kanit&display=swap"
            rel="stylesheet"
          />
        </Head>
        <body style={{ margin: 0 }}>
          <Main />
          <NextScript />
        </body>
      </Html>
    );
  }
}

export default MyDocument;
