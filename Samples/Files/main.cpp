//
// Created by devil on 18.05.17.
//

#include "engine.h"

Paranoia::Engine *engine;


int main() {
    engine = new Paranoia::Engine(ENGINE_PC);

    engine->Init();

    System::cFile File("test.txt", 0);

    File.Open(FILE_OPEN_TYPE::OPEN_READ);

    System::cFileData *fd;

    fd = File.Read(100);

    if (fd != NULL)
        std::cout << fd->data << std::endl;
    else {
        std::cout << "File not found!" << std::endl;
        engine->log->AddMessage("File not found!", LOG_TYPE::LOG_ERROR);
    }

    engine->Start();

    return 0;
}