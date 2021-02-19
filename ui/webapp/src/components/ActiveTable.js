import TableView from "./TableView";

export default function ActiveTable({ activeStatus }) {
  return <TableView data={activeStatus} listName={"Active"} />;
}
