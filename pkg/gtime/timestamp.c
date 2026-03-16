#include "timestamp.h"

long long get_system_timestamp() {
    return (long long)time(NULL);
}