import { useHistory } from "react-router-dom";
import DataTable from "react-data-table-component";

const columnify = (keynames) => {
  return keynames.map((element) => {
    return { name: element, selector: element, sortable: true };
  });
};

export default function TableView({ listName, data }) {
  const keynames = data.length > 0 ? Object.keys(data[0]) : [];
  const columns = columnify(keynames);

  const history = useHistory();

  return (
    <DataTable
      title={`${listName} List`}
      columns={columns}
      data={data}
      persistTableHead
      highlightOnHover
      onRowClicked={(row) =>
        history.push(`/${listName.toLowerCase()}/${row.uuid}`)
      }
    />
  );
}
