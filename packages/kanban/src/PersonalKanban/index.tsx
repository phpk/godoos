import React from "react";

import ThemeProvider from "PersonalKanban/providers/ThemeProvider";
import KanbanBoardContainer from "PersonalKanban/containers/KanbanBoard";
import TranslationProvider from "./providers/TranslationProvider";
import { TitleContextProvider } from "PersonalKanban/containers/KanbanBoard/title";

interface PersonalKanbanProps {}

const PersonalKanban: React.FC<PersonalKanbanProps> = () => {
  return (
    <ThemeProvider>
      <TranslationProvider>
        <TitleContextProvider>
          <KanbanBoardContainer />
        </TitleContextProvider>
      </TranslationProvider>
    </ThemeProvider>
  );
};

export default PersonalKanban;
