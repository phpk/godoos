import React, { useContext } from "react";

import AppBar from "@material-ui/core/AppBar";
import Box from "@material-ui/core/Box";
import Button from "@material-ui/core/Button";
import Dialog from "@material-ui/core/Dialog";
import TextField from "@material-ui/core/TextField";
import DialogContent from "@material-ui/core/DialogContent";
import Divider from "@material-ui/core/Divider";
import MuiToolbar from "@material-ui/core/Toolbar";
import Menu from "@material-ui/core/Menu";
import MenuItem from "@material-ui/core/MenuItem";
import Grid from "@material-ui/core/Grid";
import Typography from "@material-ui/core/Typography";
//import useMediaQuery from "@material-ui/core/useMediaQuery";
import { makeStyles } from "@material-ui/core/styles";

import { useTranslation } from "PersonalKanban/providers/TranslationProvider";
import ColumnForm from "PersonalKanban/components/ColumnForm";
import IconButton from "PersonalKanban/components/IconButton";
import { Column } from "PersonalKanban/types";
import { useTheme } from "PersonalKanban/providers/ThemeProvider";
import StorageService from "PersonalKanban/services/StorageService";
import { TitleContext } from "../KanbanBoard/title";

type AddColumnButtonProps = {
  onSubmit: any;
};

const AddColumnButton: React.FC<AddColumnButtonProps> = (props) => {
  const { onSubmit } = props;

  const { t } = useTranslation();

  const [open, setOpen] = React.useState(false);

  const handleOpenDialog = React.useCallback(() => {
    setOpen(true);
  }, []);

  const handleCloseDialog = React.useCallback(() => {
    setOpen(false);
  }, []);

  const handleSubmit = React.useCallback(
    (column: Column) => {
      onSubmit({ column });
      handleCloseDialog();
    },
    [onSubmit, handleCloseDialog]
  );

  return (
    <Box display="block">
      <IconButton icon="add" color="primary" onClick={handleOpenDialog}>
        {t("addColumn")}
      </IconButton>
      <Dialog onClose={handleCloseDialog} open={open}>
        <DialogContent>
          <ColumnForm onSubmit={handleSubmit} onCancel={handleCloseDialog} />
        </DialogContent>
      </Dialog>
    </Box>
  );
};

type ClearBoardButtonProps = {
  onClear: any;
  disabled?: boolean;
};

const ClearBoardButton: React.FC<ClearBoardButtonProps> = (props) => {
  const { disabled, onClear } = props;

  const { t } = useTranslation();

  const [open, setOpen] = React.useState(false);

  const handleOpenDialog = React.useCallback(() => {
    setOpen(true);
  }, []);

  const handleCloseDialog = React.useCallback(() => {
    setOpen(false);
  }, []);

  const handleClear = React.useCallback(
    (e) => {
      onClear({ e });
      handleCloseDialog();
    },
    [onClear, handleCloseDialog]
  );

  return (
    <Box display="flex">
      <IconButton
        icon="delete"
        color="primary"
        disabled={disabled}
        onClick={handleOpenDialog}
      ></IconButton>
      <Dialog onClose={handleCloseDialog} open={open}>
        <DialogContent>
          <Grid container spacing={1}>
            <Grid item xs={12}>
              <Typography gutterBottom variant="h6">
                {t("clearBoard")}
              </Typography>
              <Divider />
            </Grid>
            <Grid item xs={12}>
              <Typography gutterBottom>
                {t("clearBoardConfirmation")}
              </Typography>
            </Grid>
            <Grid item xs={12}>
              <Button variant="outlined" onClick={handleCloseDialog}>
                {t("cancel")}
              </Button>
              &nbsp;
              <Button color="primary" variant="contained" onClick={handleClear}>
                {t("clear")}
              </Button>
            </Grid>
          </Grid>
        </DialogContent>
      </Dialog>
    </Box>
  );
};

type LanguageButtonProps = {};

const LanguageButton: React.FC<LanguageButtonProps> = (props) => {
  const { i18n } = useTranslation();

  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleChangeLanguage = (lng: string) => () => {
    i18n.changeLanguage(lng);
    handleClose();
  };

  return (
    <Box display="block">
      <IconButton
        icon={"language"}
        aria-controls="language-menu"
        aria-haspopup="true"
        color="inherit"
        onClick={handleClick}
      />
      <Menu
        id="language-menu"
        anchorEl={anchorEl}
        keepMounted
        open={Boolean(anchorEl)}
        onClose={handleClose}
      >
        <MenuItem onClick={handleChangeLanguage("cn")}>中文</MenuItem>
        <MenuItem onClick={handleChangeLanguage("en")}>English</MenuItem>
      </Menu>
    </Box>
  );
};

const DarkThemeButton: React.FC<{}> = () => {
  const { darkTheme, handleToggleDarkTheme } = useTheme();

  return (
    <IconButton
      color="inherit"
      icon={darkTheme ? "invertColors" : "invertColorsOff"}
      onClick={handleToggleDarkTheme}
    />
  );
};

const SaveButton: React.FC<{}> = () => {
  const { title } = useContext(TitleContext);
  const SaveAction = () => {   
    const content = StorageService.getColumns();
    //console.log(dTitle)
    //console.log(content)
    const save = {
      data: JSON.stringify({ content, title }),
      type: 'exportKanban'
    }
    window.parent.postMessage(save, '*')
  }
  //const { t } = useTranslation();
  return (
    <IconButton
      color="inherit"
      icon="save"
      onClick={SaveAction}
    />
  );
};
const TitleInput: React.FC<{}> = () => {
  const { t } = useTranslation();
  //const [dataTitle, setDataTitle] = useState('未命名看板');
  // let dataTitle = localStorage.getItem("__kanbantitle")
  // dataTitle = dataTitle ? dataTitle : t("undefineTitle")
  // const [daTitle, setDaTitle] = useState(dataTitle);
  const { title, setTitle } = useContext(TitleContext); // 从上下文中获取title和setTitle
  // React.useEffect(() => {
  //   console.log("标题已更新为:", title); // 用于调试，确认title是否更新
  //   // 在这里可以执行标题更新后的其他逻辑，如DOM操作或API调用
  // }, [title]); // 监听title变化
  return (
    <Box 
    display="block" 
    component="form"
    >
        <TextField
        name="title"
        label={t("personalKanban")}
        id="dataTitle"
        variant="standard"
        value={title}
        onChange={(e:any) => {
          //localStorage.setItem("__kanbantitle", e.target.value)
          setTitle(e.target.value);
          //dataTitle = e.target.value;
        }}
      />
    </Box>
  );
};


const useToolbarStyles = makeStyles(() => ({
  paper: {
    padding: 0,
  },
}));

type ToolbarProps = {
  clearButtonDisabled?: boolean;
  onNewColumn: any;
  onClearBoard: any;
  //dataTitle?:any;
};

const Toolbar: React.FC<ToolbarProps> = (props) => {
  const { clearButtonDisabled, onNewColumn, onClearBoard } = props;

  //const { t } = useTranslation();

  const classes = useToolbarStyles();

  //const muiTheme = useMuiTheme();
  //const [dataTitle, setDataTitle] = useState('')


  //const isMobile = useMediaQuery(muiTheme.breakpoints.down("sm"));
  //dataTitle = dataTitle ? dataTitle : t("undefineTitle")
  return (
    <AppBar color="default" elevation={6} className={classes.paper}>
      <MuiToolbar>
        <Box display="flex" alignItems="center">
          <IconButton
            icon="personalKanban"
            color="primary"
            size="small"
            iconProps={{ fontSize: "large" }}
            disableRipple
            disableTouchRipple
            disableFocusRipple
          />
          &nbsp;
          <TitleInput />
        </Box>
        <Box display="flex" flexGrow={1} />
        <Box display="flex">
          <AddColumnButton onSubmit={onNewColumn} />
          &nbsp;
          <ClearBoardButton
            disabled={clearButtonDisabled}
            onClear={onClearBoard}
          />
          &nbsp;
          <LanguageButton /> &nbsp;
          <DarkThemeButton /> &nbsp;
          <SaveButton />
        </Box>
      </MuiToolbar>
    </AppBar>
  );
};

export default Toolbar;
