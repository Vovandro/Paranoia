//
// Created by devil on 18.05.17.
//

#include "engine.h"

Paranoia::Engine *engine;


int main() {
    engine = new Paranoia::Engine(ENGINE_PC);

    engine->Init("engine.cf");

    System::cFile File(engine, "log", 0);

    File.Open(FILE_OPEN_TYPE::OPEN_READ);

    System::cFileData *fd;

    fd = File.Read();

    if (fd != NULL)
        std::cout << fd->data << std::endl;
    else {
        std::cout << "File not found!" << std::endl;
        engine->log->AddMessage("File not found!", LOG_TYPE::LOG_ERROR);
    }

    engine->Start();

    delete engine;

    return 0;
}