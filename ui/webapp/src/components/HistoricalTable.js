import TableView from "./TableView";

export default function HistoricalTable({ historicalStatus }) {
  return <TableView data={historicalStatus} listName={"Historical"} />;
}
