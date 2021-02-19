import DataTable from "react-data-table-component";

const columnify = (keynames) => {
  return keynames.map((element) => {
    return { name: element, selector: element, sortable: true };
  });
};

export default function TaskView({ tasks, setFieldCallback }) {
  const taskKeynames = tasks.length > 0 ? Object.keys(tasks[0]) : [];
  const taskColumns = columnify(taskKeynames);
  return (
    <div>
      <DataTable
        title="Task List"
        columns={taskColumns}
        data={tasks}
        persistTableHead
        highlightOnHover
        onRowClicked={(row) => console.log(row)}
      />
    </div>
  );
}
