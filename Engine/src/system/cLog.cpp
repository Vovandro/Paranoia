//
// Created by devil on 19.05.17.
//

#include "../../include/system/cLog.h"
#include "../../include/engine.h"

System::cLog::cLog(Paranoia::Engine *engine, std::string fName) : System::cThread(engine, "log", 0, true, true, 10, true) {
    this->engine = engine;
    cyrMessage = NULL;
    lastMessage = NULL;

    file = new cFile(engine, fName, 0, true);
    engine->files->AddObject(file);
    file->Open(FILE_OPEN_TYPE::OPEN_WRITE_CLEAR);

    engine->threads->AddObject(this);
}

System::cLog::~cLog() {
    LockLocal();
    file->Close();
    delete file;
    UnLockLocal();
}

void System::cLog::Work() {
    LockLocal();
    if (cyrMessage == NULL) {
        UnLockLocal();
        SleepThis(100);
        return;
    }

    Write();
    UnLockLocal();
}

void System::cLog::Write() {
    if (cyrMessage == NULL)
        return;

    cLogMessage *msg = cyrMessage;
    file->Write(msg->Message + "\n");

    if (lastMessage == cyrMessage) {
        lastMessage = NULL;
    }

    cyrMessage = msg->nextMessage;

    delete msg;
}

void System::cLog::EndWork() {
}

void System::cLog::AddMessage(std::string Message, LOG_TYPE Type) {
    LockLocal();
    cLogMessage *msg = new cLogMessage();

    msg->Message = Message;
    msg->Type = Type;

    if (lastMessage == NULL) {
        cyrMessage = msg;
        lastMessage = msg;
    } else {
        lastMessage->nextMessage = msg;
        lastMessage = msg;
    }

    std::cout << Message + "\n";

    UnLockLocal();
}