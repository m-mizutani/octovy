import React from "react";
import Button from "@mui/material/Button";
import Alert from "@mui/material/Alert";
import Grid from "@mui/material/Grid";
import Fade from "@mui/material/Fade";

export default function RawData(props: { scanID: string }) {
  const [state, setState] = React.useState<{
    err?: any;
    data?: any;
    isLoaded?: boolean;
  }>({});
  const [copied, setCopied] = React.useState<boolean>(false);

  const get = function () {
    fetch(`/api/v1/scan/${props.scanID}/report`)
      .then((res) => res.json())
      .then(
        (result) => {
          console.log("result:", { result });
          if (result.error) {
            setState({ isLoaded: true, err: result.error });
          } else {
            setState({
              isLoaded: true,
              data: result.data,
            });
          }
        },
        (error) => {
          console.log("error:", { error });
          setState({
            isLoaded: true,
            err: error,
          });
        }
      );
  };

  React.useEffect(get, []);
  React.useEffect(() => {
    if (!copied) {
      return;
    }

    const timer = setTimeout(() => {
      setCopied(false);
    }, 3 * 1000);

    return () => {
      clearTimeout(timer);
    };
  }, [copied]);

  if (!state.isLoaded) {
    return <Alert severity="info">Loading</Alert>;
  }
  if (state.err) {
    return <Alert severity="error">{state.err}</Alert>;
  }

  return (
    <Grid container spacing={2}>
      <Grid item>
        <Button
          variant="outlined"
          size="small"
          onClick={() => {
            navigator.clipboard.writeText(JSON.stringify(state.data)).then(
              () => {
                setCopied(true);
              },
              () => {}
            );
          }}>
          Copy report
        </Button>
      </Grid>
      <Grid item>
        <Fade in={copied}>
          <Alert severity="success">Copied</Alert>
        </Fade>
      </Grid>
    </Grid>
  );
}
