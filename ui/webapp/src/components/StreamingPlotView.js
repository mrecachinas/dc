import React from "react";
import { SigPlot, WPipeLayer } from "react-sigplot";

export default function StreamingPlotView({ streamID }) {
  // TODO: Get base URL...
  const wsurl = `/${streamID}`;
  return (
    <div>
      <SigPlot>
        <WPipeLayer wsurl={wsurl} />
      </SigPlot>
    </div>
  );
}
