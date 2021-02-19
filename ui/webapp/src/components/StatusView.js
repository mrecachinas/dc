import ActiveTable from "./ActiveTable";
import HistoricalTable from "./HistoricalTable";

export default function StatusView({ status }) {
  const { active: activeStatus, historical: historicalStatus } = status;

  return (
    <div>
      <ActiveTable activeStatus={activeStatus} />
      <HistoricalTable historicalStatus={historicalStatus} />
    </div>
  );
}
