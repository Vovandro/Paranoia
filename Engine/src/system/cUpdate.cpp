//
// Created by devil on 25.05.17.
//

#include "../../include/system/cUpdate.h"
#include "../../include/engine.h"

System::cUpdate::cUpdate(Paranoia::Engine *engine): cThread(engine->threads, "update", 1, true, true, 1, true) {

}