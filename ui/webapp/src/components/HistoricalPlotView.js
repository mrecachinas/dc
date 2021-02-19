import React from "react";
import { SigPlot, HrefLayer } from "react-sigplot";

export default function HistoricalPlotView({ historicalID }) {
  // TODO: Get base URL...
  // const href = `/api/${historicalID}`;
  const href = "https://sigplot.lgsinnovations.com/dat/penny.prm";
  return (
    <div>
      <SigPlot height={400} width={"100%"}>
        <HrefLayer href={href} />
      </SigPlot>
    </div>
  );
}
