//
// Created by devil on 25.05.17.
//

#ifndef PROJECT_CSTATEMANAGER_H
#define PROJECT_CSTATEMANAGER_H

#include "cState.h"

namespace Core {
    class cStateManager {
    protected:
        cState *state;
    public:
        cStateManager();
        ~cStateManager();

        //Проталкиваем новое состояние
        void Push(cState *newState);
        //Снимаем последнее
        void Pop();
        //Удаляем все состояния. Если isMessage = false то не уведомляем о удалении
        void PopAll(bool isMessage = true);
        //Получение текущего
        cState * Get();

        //Обновление
        void Update(int dt);
    };
}

#endif //PROJECT_CSTATEMANAGER_H
