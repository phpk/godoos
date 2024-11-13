import React, { useContext } from "react";

import Box from "@material-ui/core/Box";
import { makeStyles } from "@material-ui/core/styles";

import KanbanBoard from "PersonalKanban/components/KanbanBoard";
import { Column, Record } from "PersonalKanban/types";
import {
  getId,
  getCreatedAt,
  getInitialState,
  reorder,
  reorderCards,
} from "PersonalKanban/services/Utils";
import StorageService from "PersonalKanban/services/StorageService";
import Toolbar from "PersonalKanban/containers/Toolbar";
import { TitleContext } from "./title";

const useKanbanBoardStyles = makeStyles((theme) => ({
  toolbar: theme.mixins.toolbar,
}));

type KanbanBoardContainerProps = {};

//let initialState = StorageService.getColumns();

//if (!initialState) {
let initialState = getInitialState();
// localStorage.setItem("__kanbantitle", "");
//}
function isBase64(str : string) {
  if (str === '' || str.trim() === '') {
    return false
  }
  try {
    return btoa(atob(str)) === str
  } catch (err) {
    return false
  }
}
function decodeBase64(base64String : string) {
  // 将Base64字符串分成每64个字符一组
  const padding =
    base64String.length % 4 === 0 ? 0 : 4 - (base64String.length % 4)
  base64String += '='.repeat(padding)

  // 使用atob()函数解码Base64字符串
  const binaryString = atob(base64String)

  // 将二进制字符串转换为TypedArray
  const bytes = new Uint8Array(binaryString.length)
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i)
  }

  // 将TypedArray转换为字符串
  return new TextDecoder('utf-8').decode(bytes)
}
const KanbanBoardContainer: React.FC<KanbanBoardContainerProps> = (props) => {
  const [columns, setColumns] = React.useState<Column[]>(initialState);
  const { setTitle } = useContext(TitleContext); // 从上下文中获取setTitle函数
  const EventHandler = (e: any) => {
    const eventData = e.data
    //const titleInput = document.querySelector<HTMLInputElement>("#dataTitle");
    // console.log(eventData)
    if (eventData.type === 'start') {
      if (eventData.title) {
        const baseTitle = eventData.title.substring(0, eventData.title.lastIndexOf('.'))
        // localStorage.setItem("__kanbantitle", baseTitle)
        // setTimeout(() => {
        //   titleInput!.value = baseTitle
        // }, 0)
        setTitle(baseTitle)
      }
    }
    if (eventData.type === 'init') {
      const data = eventData.data
      if (data) {
        if(data.content) {
          if (typeof data.content === 'string' && isBase64(data.content)) {
            data.content = decodeBase64(data.content)
            data.content = JSON.parse(data.content)
          }
          initialState = data.content;
          setColumns(data.content);
        }
        
        if (data.title) {
          const baseTitle = data.title.substring(0, data.title.lastIndexOf('.'))
          // localStorage.setItem("__kanbantitle", baseTitle)
          // setTimeout(() => {
          //   console.log("Setting title to:", baseTitle);
          //   setTitle(baseTitle);
          //   console.log("setTitle function has been called.");
          // }, 100)
          setTitle(baseTitle);
          
        }
        
      }
    }
  }
  window.addEventListener('load', () => {
    window.parent.postMessage({ type: 'initSuccess' }, '*')
    window.addEventListener('message', EventHandler)
  })
  window.addEventListener('unload', () => {
    window.removeEventListener('message', EventHandler)
  })
  const classes = useKanbanBoardStyles();

  const cloneColumns = React.useCallback((columns: Column[]) => {
    return columns.map((column: Column) => ({
      ...column,
      records: [...column.records!],
    }));
  }, []);

  const getColumnIndex = React.useCallback(
    (id: string) => {
      return columns.findIndex((c: Column) => c.id === id);
    },
    [columns]
  );

  const getRecordIndex = React.useCallback(
    (recordId: string, columnId: string) => {
      return columns[getColumnIndex(columnId)]?.records?.findIndex(
        (r: Record) => r.id === recordId
      );
    },
    [columns, getColumnIndex]
  );

  const handleClearBoard = React.useCallback(() => {
    setColumns([]);
  }, []);

  const handleAddColumn = React.useCallback(
    ({ column }: { column: Column }) => {
      setColumns((columns: Column[]) => [
        ...columns,
        Object.assign(
          { id: getId(), records: [], createdAt: getCreatedAt() },
          column
        ),
      ]);
    },
    []
  );

  const handleColumnMove = React.useCallback(
    ({ column, index }: { column: Column; index: number }) => {
      const updatedColumns = reorder(columns, getColumnIndex(column.id), index);
      setColumns(updatedColumns);
    },
    [columns, getColumnIndex]
  );

  const handleColumnEdit = React.useCallback(
    ({ column }: { column: Column }) => {
      setColumns((_columns: Column[]) => {
        const columnIndex = getColumnIndex(column.id);
        const columns = cloneColumns(_columns);
        columns[columnIndex].title = column.title;
        columns[columnIndex].description = column.description;
        columns[columnIndex].color = column.color;
        columns[columnIndex].wipEnabled = column.wipEnabled;
        columns[columnIndex].wipLimit = column.wipLimit;
        return columns;
      });
    },
    [getColumnIndex, cloneColumns]
  );

  const handleColumnDelete = React.useCallback(
    ({ column }: { column: Column }) => {
      setColumns((_columns: Column[]) => {
        const columns = cloneColumns(_columns);
        columns.splice(getColumnIndex(column.id), 1);
        return columns;
      });
    },
    [cloneColumns, getColumnIndex]
  );

  const handleCardMove = React.useCallback(
    ({
      column,
      index,
      source,
      record,
    }: {
      column: Column;
      index: number;
      source: Column;
      record: Record;
    }) => {
      const updatedColumns = reorderCards({
        columns,
        destinationColumn: column,
        destinationIndex: index,
        sourceColumn: source,
        sourceIndex: getRecordIndex(record.id, source.id)!,
      });

      setColumns(updatedColumns);
    },
    [columns, getRecordIndex]
  );

  const handleAddRecord = React.useCallback(
    ({ column, record }: { column: Column; record: Record }) => {
      const columnIndex = getColumnIndex(column.id);
      setColumns((_columns: Column[]) => {
        const columns = cloneColumns(_columns);

        columns[columnIndex].records = [
          {
            id: getId(),
            title: record.title,
            description: record.description,
            color: record.color,
            createdAt: getCreatedAt(),
          },
          ...columns[columnIndex].records,
        ];
        return columns;
      });
    },
    [cloneColumns, getColumnIndex]
  );

  const handleRecordEdit = React.useCallback(
    ({ column, record }: { column: Column; record: Record }) => {
      const columnIndex = getColumnIndex(column.id);
      const recordIndex = getRecordIndex(record.id, column.id);
      setColumns((_columns) => {
        const columns = cloneColumns(_columns);
        const _record = columns[columnIndex].records[recordIndex!];
        _record.title = record.title;
        _record.description = record.description;
        _record.color = record.color;
        return columns;
      });
    },
    [getColumnIndex, getRecordIndex, cloneColumns]
  );

  const handleRecordDelete = React.useCallback(
    ({ column, record }: { column: Column; record: Record }) => {
      const columnIndex = getColumnIndex(column.id);
      const recordIndex = getRecordIndex(record.id, column.id);
      setColumns((_columns) => {
        const columns = cloneColumns(_columns);
        columns[columnIndex].records.splice(recordIndex!, 1);
        return columns;
      });
    },
    [cloneColumns, getColumnIndex, getRecordIndex]
  );

  const handleAllRecordDelete = React.useCallback(
    ({ column }: { column: Column }) => {
      const columnIndex = getColumnIndex(column.id);
      setColumns((_columns) => {
        const columns = cloneColumns(_columns);
        columns[columnIndex].records = [];
        return columns;
      });
    },
    [cloneColumns, getColumnIndex]
  );

  React.useEffect(() => {
    StorageService.setColumns(columns);
  }, [columns]);

  return (
    <>  
      <Toolbar
        clearButtonDisabled={!columns.length}
        onNewColumn={handleAddColumn}
        onClearBoard={handleClearBoard}
      />
      <div className={classes.toolbar} />
      <Box padding={1}>
        <KanbanBoard
          columns={columns}
          onColumnMove={handleColumnMove}
          onColumnEdit={handleColumnEdit}
          onColumnDelete={handleColumnDelete}
          onCardMove={handleCardMove}
          onAddRecord={handleAddRecord}
          onRecordEdit={handleRecordEdit}
          onRecordDelete={handleRecordDelete}
          onAllRecordDelete={handleAllRecordDelete}
        />
      </Box> 
    </>
  );
};

export default KanbanBoardContainer;
