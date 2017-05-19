//
// Created by devil on 19.05.17.
//

#ifndef PROJECT_CLOG_H
#define PROJECT_CLOG_H

#include "cThread.h"
#include "cFile.h"

enum LOG_TYPE {
    LOG_DEBUG,
    LOG_MESSAGE,
    LOG_WARNING,
    LOG_ERROR,
    LOG_CRITICAL,
};

namespace Paranoia {
    class Engine;
}

namespace System {
    class cLogMessage {
    public:
        cLogMessage() {nextMessage = NULL;};

        std::string Message;
        LOG_TYPE Type;
        cLogMessage *nextMessage;
    };

    class cLog : public cThread {
    protected:
        Paranoia::Engine *engine;
        cFile *file;

        cLogMessage *cyrMessage;
        cLogMessage *lastMessage;

        //Запись текущего сообщения
        void Write();

    public:
        cLog(Paranoia::Engine *engine, std::string fName);
        ~cLog();

        //Действие выполняемое в потоке
        virtual void Work() override;

        //Действие выполняемое при закрытие потока
        virtual void EndWork() override;

        //Добавление сообщеня в очередь записи
        void AddMessage(std::string Message, LOG_TYPE Type);
    };
}

#endif //PROJECT_CLOG_H
