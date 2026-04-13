export const RefreshTable = "props-table:refresh";
export type RefreshTableData<T> = T[];
export type RefresTablehHandler<T> = (
  event: CustomEvent<RefreshTableData<T>>,
) => void;
export type RefreshTableEvent<T> = CustomEvent<RefreshTableData<T>>;

export type ServerItem = {
  id: string;
  status: string;
  stats: {
    cpu: {
      usage: number;
    };
    mem: {
      usage: number;
    };
  };
};
