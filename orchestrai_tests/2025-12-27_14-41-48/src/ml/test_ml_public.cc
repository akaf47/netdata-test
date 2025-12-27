// SPDX-License-Identifier: GPL-3.0-or-later

#include <gtest/gtest.h>
#include <gmock/gmock.h>
#include <cstring>
#include <vector>
#include <unordered_map>
#include <atomic>

// Mock structures and types
struct RRDHOST {
    char machine_guid[50];
    char *hostname;
    void *ml_host;
    int dummy;
};

struct RRDSET {
    RRDHOST *rrdhost;
    char *id;
    char *name;
    void *ml_chart;
    int dummy;
};

struct RRDDIM {
    RRDSET *rrdset;
    char *id;
    char *name;
    void *ml_dimension;
    int dummy;
};

struct BUFFER {
    std::string content;
    int dummy;
};

enum METRIC_TYPE {
    METRIC_TYPE_CONSTANT = 0,
    METRIC_TYPE_OTHER = 1
};

enum TRAINING_STATUS {
    TRAINING_STATUS_UNTRAINED = 0,
    TRAINING_STATUS_TRAINED = 1,
    TRAINING_STATUS_PENDING_WITH_MODEL = 2,
    TRAINING_STATUS_PENDING_WITHOUT_MODEL = 3,
    TRAINING_STATUS_SILENCED = 4
};

enum MACHINE_LEARNING_STATUS {
    MACHINE_LEARNING_STATUS_ENABLED = 0,
    MACHINE_LEARNING_STATUS_DISABLED_DUE_TO_EXCLUDED_CHART = 1
};

enum ML_QUEUE_ITEM_TYPE {
    ML_QUEUE_ITEM_TYPE_CREATE_NEW_MODEL = 0
};

typedef double calculated_number_t;
typedef void (*spinlock_t);
typedef void (*mutex_t);
typedef void (*queue_t);
typedef void (*nd_thread_t);
typedef unsigned long time_t;

struct ml_machine_learning_stats_t {
    uint64_t num_anomalous_dimensions = 0;
    uint64_t num_normal_dimensions = 0;
    uint64_t num_training_status_trained = 0;
    uint64_t num_training_status_pending_with_model = 0;
    uint64_t num_training_status_untrained = 0;
    uint64_t num_training_status_pending_without_model = 0;
    uint64_t num_training_status_silenced = 0;
};

struct ml_kmeans_t {
    int dummy;
};

struct ml_dimension_t {
    RRDDIM *rd;
    METRIC_TYPE mt;
    TRAINING_STATUS ts;
    uint64_t suppression_anomaly_counter;
    uint64_t suppression_window_counter;
    bool training_in_progress;
    ml_kmeans_t kmeans;
    std::vector<void*> cns;
    spinlock_t slock;
    std::vector<void*> km_contexts;
    MACHINE_LEARNING_STATUS mls;
};

struct ml_chart_t {
    RRDSET *rs;
    ml_machine_learning_stats_t mls;
};

struct ml_metrics_statistics {
    uint64_t anomalous;
    uint64_t normal;
    uint64_t trained;
    uint64_t pending;
    uint64_t silenced;
};

struct ml_worker_t {
    size_t id;
    queue_t queue;
    calculated_number_t *training_cns;
    calculated_number_t *scratch_training_cns;
    nd_thread_t nd_thread;
};

struct DimensionLookupInfo {
    const char *host_guid;
    const char *chart_id;
    const char *dimension_id;
    
    DimensionLookupInfo(const char *hg, const char *cid, const char *did)
        : host_guid(hg), chart_id(cid), dimension_id(did) {}
};

struct ml_request_create_new_model_t {
    DimensionLookupInfo DLI;
};

struct ml_queue_item_t {
    ML_QUEUE_ITEM_TYPE type;
    ml_request_create_new_model_t create_new_model;
};

struct ml_config_t {
    bool enable_anomaly_detection;
    bool stream_anomaly_detection_charts;
    void *sp_host_to_skip;
    void *sp_charts_to_skip;
    uint64_t training_window;
    uint64_t min_training_window;
    uint64_t max_training_vectors;
    uint64_t max_samples_to_smooth;
    uint64_t train_every;
    uint64_t diff_n;
    uint64_t lag_n;
    uint64_t max_kmeans_iters;
    double dimension_anomaly_score_threshold;
    int anomaly_detection_grouping_method;
    int64_t anomaly_detection_query_duration;
    std::string hosts_to_skip;
    std::string charts_to_skip;
    std::vector<uint32_t> random_nums;
    std::vector<ml_worker_t> workers;
    size_t num_worker_threads;
    size_t flush_models_batch_size;
    bool detection_stop;
    bool training_stop;
    nd_thread_t detection_thread;
};

struct ml_host_t {
    RRDHOST *rh;
    ml_machine_learning_stats_t mls;
    double host_anomaly_rate;
    void *anomaly_rate_rs;
    queue_t queue;
    mutex_t mutex;
    spinlock_t context_anomaly_rate_spinlock;
    bool ml_running;
    std::atomic<bool> reset_pointers;
};

typedef void *rrd_ml_host_t;
typedef void *rrd_ml_chart_t;
typedef void *rrd_ml_dimension_t;

// Mock global config
ml_config_t Cfg;
sqlite3 *ml_db = nullptr;

// Mock functions
void netdata_mutex_init(mutex_t *) {}
void netdata_mutex_destroy(mutex_t *) {}
void netdata_mutex_lock(mutex_t *) {}
void netdata_mutex_unlock(mutex_t *) {}
void spinlock_init(spinlock_t *) {}
void spinlock_lock(spinlock_t *) {}
void spinlock_unlock(spinlock_t *) {}
void sleep_usec(unsigned long) {}
void netdata_log_error(const char *, ...) {}

bool simple_pattern_matches(void *, const char *) { return false; }
const char *rrdhost_hostname(RRDHOST *rh) { return rh->hostname; }
const char *rrdset_name(RRDSET *rs) { return rs->name; }
const char *rrdset_id(RRDSET *rs) { return rs->id; }
const char *rrddim_id(RRDDIM *rd) { return rd->id; }

void rrdset_foreach_read(void *&ptr, RRDHOST *rh) {}
void rrdset_foreach_done(void *ptr) {}
void rrddim_foreach_read(void *&ptr, RRDSET *rs) {}
void rrddim_foreach_done(void *ptr) {}

void ml_config_load(ml_config_t *cfg) {}
void ml_kmeans_init(ml_kmeans_t *) {}
void metaqueue_ml_load_models(RRDDIM *) {}
void ml_queue_push(queue_t, ml_queue_item_t) {}
queue_t ml_queue_init() { return nullptr; }
void ml_queue_destroy(queue_t) {}
void ml_queue_signal(queue_t) {}

void error_report(const char *, ...) {}
void global_statistics_ml_models_deserialization_failures() {}
bool ml_dimension_deserialize_kmeans(const char *) { return true; }

int sqlite3_open(const char *, sqlite3 **) { return 0; }
int sqlite3_exec(sqlite3 *, const char *, int (*)(void*, int, char**, char**), void *, char **) { return 0; }
int sqlite3_close(sqlite3 *) { return 0; }
const char *sqlite3_errstr(int) { return ""; }
void sqlite3_free(void *) {}

int perform_ml_database_migration(sqlite3 *, int) { return 0; }
bool configure_sqlite_database(sqlite3 *, int, const char *) { return false; }
uint64_t sqlite_get_db_space(sqlite3 *) { return 0; }
void sql_close_database(sqlite3 *, const char *) {}

nd_thread_t nd_thread_create(const char *, int, void (*)(void *), void *) { return nullptr; }
void nd_thread_join(nd_thread_t) {}

const char *time_grouping_id2txt(int) { return "test"; }
bool ml_dimension_predict(ml_dimension_t *, double, bool) { return false; }
void ml_chart_update_dimension(ml_chart_t *, ml_dimension_t *, bool) {}

const char *db_models_create_table = "CREATE TABLE IF NOT EXISTS test (id INTEGER);";
char netdata_configured_cache_dir[] = "/tmp";

#include "ml_public.cc"

// Test class
class MLPublicTest : public ::testing::Test {
protected:
    void SetUp() override {
        Cfg.enable_anomaly_detection = true;
        Cfg.stream_anomaly_detection_charts = false;
        Cfg.sp_host_to_skip = nullptr;
        Cfg.sp_charts_to_skip = nullptr;
        Cfg.training_window = 3600;
        Cfg.min_training_window = 900;
        Cfg.max_training_vectors = 10000;
        Cfg.max_samples_to_smooth = 100;
        Cfg.train_every = 600;
        Cfg.diff_n = 1;
        Cfg.lag_n = 5;
        Cfg.max_kmeans_iters = 1000;
        Cfg.dimension_anomaly_score_threshold = 0.5;
        Cfg.anomaly_detection_grouping_method = 1;
        Cfg.anomaly_detection_query_duration = 3600000;
        Cfg.hosts_to_skip = "";
        Cfg.charts_to_skip = "";
        Cfg.num_worker_threads = 1;
        Cfg.flush_models_batch_size = 10;
        Cfg.detection_stop = false;
        Cfg.training_stop = false;
        Cfg.detection_thread = nullptr;
        Cfg.workers.clear();
        Cfg.random_nums.clear();
        ml_db = nullptr;
    }

    void TearDown() override {
        Cfg.workers.clear();
        Cfg.random_nums.clear();
    }

    RRDHOST *createTestHost() {
        RRDHOST *rh = new RRDHOST();
        strcpy(rh->machine_guid, "test-guid-12345678");
        rh->hostname = strdup("test-host");
        rh->ml_host = nullptr;
        return rh;
    }

    RRDSET *createTestChart(RRDHOST *rh) {
        RRDSET *rs = new RRDSET();
        rs->rrdhost = rh;
        rs->id = strdup("system.cpu");
        rs->name = strdup("system.cpu");
        rs->ml_chart = nullptr;
        return rs;
    }

    RRDDIM *createTestDimension(RRDSET *rs) {
        RRDDIM *rd = new RRDDIM();
        rd->rrdset = rs;
        rd->id = strdup("user");
        rd->name = strdup("user");
        rd->ml_dimension = nullptr;
        return rd;
    }

    void freeHost(RRDHOST *rh) {
        if (rh) {
            free(rh->hostname);
            delete rh;
        }
    }

    void freeChart(RRDSET *rs) {
        if (rs) {
            free(rs->id);
            free(rs->name);
            delete rs;
        }
    }

    void freeDimension(RRDDIM *rd) {
        if (rd) {
            free(rd->id);
            free(rd->name);
            delete rd;
        }
    }
};

// Tests for ml_capable()
TEST_F(MLPublicTest, ml_capable_returns_true) {
    EXPECT_TRUE(ml_capable());
}

// Tests for ml_enabled()
TEST_F(MLPublicTest, ml_enabled_returns_false_when_rh_is_null) {
    EXPECT_FALSE(ml_enabled(nullptr));
}

TEST_F(MLPublicTest, ml_enabled_returns_false_when_enable_anomaly_detection_false) {
    RRDHOST *rh = createTestHost();
    Cfg.enable_anomaly_detection = false;
    
    EXPECT_FALSE(ml_enabled(rh));
    
    freeHost(rh);
}

TEST_F(MLPublicTest, ml_enabled_returns_true_when_all_conditions_met) {
    RRDHOST *rh = createTestHost();
    Cfg.enable_anomaly_detection = true;
    
    EXPECT_TRUE(ml_enabled(rh));
    
    freeHost(rh);
}

// Tests for ml_streaming_enabled()
TEST_F(MLPublicTest, ml_streaming_enabled_returns_config_value) {
    Cfg.stream_anomaly_detection_charts = true;
    EXPECT_TRUE(ml_streaming_enabled());
    
    Cfg.stream_anomaly_detection_charts = false;
    EXPECT_FALSE(ml_streaming_enabled());
}

// Tests for ml_host_new()
TEST_F(MLPublicTest, ml_host_new_creates_host_when_enabled) {
    RRDHOST *rh = createTestHost();
    Cfg.enable_anomaly_detection = true;
    Cfg.num_worker_threads = 1;
    Cfg.workers.resize(1);
    
    ml_host_new(rh);
    
    ASSERT_NE(rh->ml_host, nullptr);
    ml_host_t *host = (ml_host_t *)rh->ml_host;
    EXPECT_EQ(host->rh, rh);
    EXPECT_FALSE(host->ml_running);
    
    ml_host_delete(rh);
    freeHost(rh);
}

TEST_F(MLPublicTest, ml_host_new_skips_when_disabled) {
    RRDHOST *rh = createTestHost();
    Cfg.enable_anomaly_detection = false;
    
    ml_host_new(rh);
    
    EXPECT_EQ(rh->ml_host, nullptr);
    
    freeHost(rh);
}

// Tests for ml_host_delete()
TEST_F(MLPublicTest, ml_host_delete_frees_resources) {
    RRDHOST *rh = createTestHost();
    Cfg.enable_anomaly_detection = true;
    Cfg.num_worker_threads = 1;
    Cfg.workers.resize(1);
    
    ml_host_new(rh);
    ml_host_delete(rh);
    
    EXPECT_EQ(rh->ml_host, nullptr);
    
    freeHost(rh);
}

TEST_F(MLPublicTest, ml_host_delete_handles_null_host) {
    RRDHOST *rh = createTestHost();
    
    ml_host_delete(rh);
    
    EXPECT_EQ(rh->ml_host, nullptr);
    
    freeHost(rh);
}

// Tests for ml_host_start()
TEST_F(MLPublicTest, ml_host_start_sets_running_flag) {
    RRDHOST *rh = createTestHost();
    Cfg.enable_anomaly_detection = true;
    Cfg.num_worker_threads = 1;
    Cfg.workers.resize(1);
    
    ml_host_new(rh);
    ml_host_start(rh);
    
    ml_host_t *host = (ml_host_t *)rh->ml_host;
    EXPECT_TRUE(host->ml_running);
    
    ml_host_delete(rh);
    freeHost(rh);
}

TEST_F(MLPublicTest, ml_host_start_handles_null_host) {
    RRDHOST *rh = createTestHost();
    
    ml_host_start(rh);
    
    EXPECT_EQ(rh->ml_host, nullptr);
    
    freeHost(rh);
}

// Tests for ml_host_stop()
TEST_F(MLPublicTest, ml_host_stop_clears_running_flag) {
    RRDHOST *rh = createTestHost();
    Cfg.enable_anomaly_detection = true;
    Cfg.num_worker_threads = 1;
    Cfg.workers.resize(1);
    
    ml_host_new(rh);
    ml_host_start(rh);
    ml_host_stop(rh);
    
    ml_host_t *host = (ml_host_t *)rh->ml_host;
    EXPECT_FALSE(host->ml_running);
    
    ml_host_delete(rh);
    freeHost(rh);
}

TEST_F(MLPublicTest, ml_host_stop_handles_null_host) {
    RRDHOST *rh = createTestHost();
    
    ml_host_stop(rh);
    
    EXPECT_EQ(rh->ml_host, nullptr);
    
    freeHost(rh);
}

TEST_F(MLPublicTest, ml_host_stop_handles_not_running_host) {
    RRDHOST *rh = createTestHost();
    Cfg.enable_anomaly_detection = true;
    Cfg.num_worker_threads = 1;
    Cfg.workers.resize(1);
    
    ml_host_new(rh);
    ml_host_stop(rh);
    
    ml_host_t *host = (ml_host_t *)rh->ml_host;
    EXPECT_FALSE(host->ml_running);
    
    ml_host_delete(rh);
    freeHost(rh);
}

// Tests for ml_host_get_info()
TEST_F(MLPublicTest, ml_host_get_info_with_null_host) {
    RRDHOST *rh = createTestHost();
    BUFFER wb;
    
    ml_host_get_info(rh, &wb);
    
    freeHost(rh);
}

TEST_F(MLPublicTest, ml_host_get_info_with_valid_host) {
    RRDHOST *rh = createTestHost();
    BUFFER wb;
    Cfg.enable_anomaly_detection = true;
    Cfg.num_worker_threads = 1;
    Cfg.workers.resize(1);
    
    ml_host_new(rh);
    ml_host_get_info(rh, &wb);
    
    ml_host_delete(rh);
    freeHost(rh);
}

// Tests for ml_host_get_detection_info()
TEST_F(MLPublicTest, ml_host_get_detection_info_with_null_host) {
    RRDHOST *rh = createTestHost();
    BUFFER wb;
    
    ml_host_get_detection_info(rh, &wb);
    
    freeHost(rh);
}

TEST_F(MLPublicTest, ml_host_get_detection_info_with_valid_host) {
    RRDHOST *rh = createTestHost();
    BUFFER wb;
    Cfg.enable_anomaly_detection = true;
    Cfg.num_worker_threads = 1;
    Cfg.workers.resize(1);
    
    ml_host_new(rh);
    ml_host_start(rh);
    ml_host_get_detection_info(rh, &wb);
    
    ml_host_delete(rh);
    freeHost(rh);
}

// Tests for ml_host_get_host_status()
TEST_F(MLPublicTest, ml_host_get_host_status_with_null_host) {
    RRDHOST *rh = createTestHost();
    ml_metrics_statistics mlm;
    
    bool result = ml_host_get_host_status(rh, &mlm);
    
    EXPECT_FALSE(result);
    EXPECT_EQ(mlm.anomalous, 0);
    EXPECT_EQ(mlm.normal, 0);
    
    freeHost(rh);
}

TEST_F(MLPublicTest, ml_host_get_host_status_with_valid_host) {
    RRDHOST *rh = createTestHost();
    ml_metrics_statistics mlm;
    Cfg.enable_anomaly_detection = true;
    Cfg.num_worker_threads = 1;
    Cfg.workers.resize(1);
    
    ml_host_new(rh);
    bool result = ml_host_get_host_status(rh, &mlm);
    
    EXPECT_TRUE(result);
    
    ml_host_delete(rh);
    freeHost(rh);
}

// Tests for ml_host_running()
TEST_F(MLPublicTest, ml_host_running_returns_false_when_null) {
    RRDHOST *rh = createTestHost();
    
    bool result = ml_host_running(rh);
    
    EXPECT_FALSE(result);
    
    freeHost(rh);
}

TEST_F(MLPublicTest, ml_host_running_returns_true_when_host_exists) {
    RRDHOST *rh = createTestHost();
    Cfg.enable_anomaly_detection = true;
    Cfg.num_worker_threads = 1;
    Cfg.workers.resize(1);
    
    ml_host_new(rh);
    bool result = ml_host_running(rh);
    
    EXPECT_TRUE(result);
    
    ml_host_delete(rh);
    freeHost(rh);
}

// Tests for ml_host_get_models()
TEST_F(MLPublicTest, ml_host_get_models_is_stub) {
    RRDHOST *rh = createTestHost();
    BUFFER wb;
    
    ml_host_get_models(rh, &wb);
    
    freeHost(rh);
}

// Tests for ml_chart_new()
TEST_F(MLPublicTest, ml_chart_new_creates_chart) {
    RRDHOST *rh = createTestHost();
    RRDSET *rs = createTestChart(rh);
    Cfg.enable_anomaly_detection = true;
    Cfg.num_worker_threads = 1;
    Cfg.workers.resize(1);
    
    ml_host_new(rh);
    ml_chart_new(rs);
    
    ASSERT_NE(rs->ml_chart, nullptr);
    ml_chart_t *chart = (ml_chart_t *)rs->ml_chart;
    EXPECT_EQ(chart->rs, rs);
    
    ml_chart_delete(rs);
    ml_host_delete(rh);
    freeChart(rs);
    freeHost(rh);
}

TEST_F(MLPublicTest, ml_chart_new_skips_when_no_host) {
    RRDHOST *rh = createTestHost();
    RRDSET *rs = createTestChart(rh);
    
    ml_chart_new(rs);
    
    EXPECT_EQ(rs->ml_chart, nullptr);
    
    freeChart(rs);
    freeHost(rh);
}

// Tests for ml_chart_delete()
TEST_F(MLPublicTest, ml_chart_delete_frees_resources) {
    RRDHOST *rh = createTestHost();
    RRDSET *rs = createTestChart(rh);
    Cfg.enable_anomaly_detection = true;
    Cfg.num_worker_threads = 1;
    Cfg.workers.resize(1);
    
    ml_host_new(rh);
    ml_chart_new(rs);
    ml_chart_delete(rs);
    
    EXPECT_EQ(rs->ml_chart, nullptr);
    
    ml_host_delete(rh);
    freeChart(rs);
    freeHost(rh);
}

// Tests for ml_chart_update_begin()
TEST_F(MLPublicTest, ml_chart_update_begin_returns_false_when_no_chart) {
    RRDHOST *rh = createTestHost();
    RRDSET *rs = createTestChart(rh);
    
    bool result = ml_chart_update_begin(rs);
    
    EXPECT_FALSE(result);
    
    freeChart(rs);
    freeHost(rh);
}

TEST_F(MLPublicTest, ml_chart_update_begin_returns_true_when_chart_exists) {
    RRDHOST *rh = createTestHost();
    RRDSET *rs = createTestChart(rh);
    Cfg.enable_anomaly_detection = true;
    Cfg.num_worker_threads = 1;
    Cfg.workers.resize(1);
    
    ml_host_new(rh);
    ml_chart_new(rs);
    bool result = ml_chart_update_begin(rs);
    
    EXPECT_TRUE(result);
    
    ml_chart_delete(rs);
    ml_host_delete(rh);
    freeChart(rs);
    freeHost(rh);
}

// Tests for ml_chart_update_end()
TEST_F(MLPublicTest, ml_chart_update_end_handles_null_chart) {
    RRDHOST *rh = createTestHost();
    RRDSET *rs = createTestChart(rh);
    
    ml_chart_update_end(rs);
    
    freeChart(rs);
    freeHost(rh);
}

TEST_F(MLPublicTest, ml_chart_update_end_with_valid_chart) {
    RRDHOST *rh = createTestHost();
    RRDSET *rs = createTestChart(rh);
    Cfg.enable_anomaly_detection = true;
    Cfg.num_worker_threads = 1;
    Cfg.workers.resize(1);
    
    ml_host_new(rh);
    ml_chart_new(rs);
    ml_chart_update_end(rs);
    
    ml_chart_delete(rs);
    ml_host_delete(rh);
    freeChart(rs);
    freeHost(rh);
}

// Tests for ml_dimension_new()
TEST_F(MLPublicTest, ml_dimension_new_creates_dimension) {
    RRDHOST *rh = createTestHost();
    RRDSET *rs = createTestChart(rh);
    RRDDIM *rd = createTestDimension(rs);
    Cfg.enable_anomaly_detection = true;
    Cfg.num_worker_threads = 1;
    Cfg.workers.resize(1);
    Cfg.num_models_to_use = 1;
    
    ml_host_new(rh);
    ml_chart_new(rs);
    ml_dimension_new(rd);
    
    ASSERT_NE(rd->ml_dimension, nullptr);
    ml_dimension_t *dim = (ml_dimension_t *)rd->ml_dimension;
    EXPECT_EQ(dim->rd, rd);
    EXPECT_EQ(dim->mt, METRIC_TYPE_CONSTANT);
    EXPECT_EQ(dim->ts, TRAINING_STATUS_UNTRAINED);
    
    ml_dimension_delete(rd);
    ml_chart_delete(rs);
    ml_host_delete(rh);
    freeDimension(rd);
    freeChart(rs);
    freeHost(rh);
}

TEST_F(MLPublicTest, ml_dimension_new_skips_when_no_chart) {
    RRDHOST *rh = createTestHost();
    RRDSET *rs = createTestChart(rh);
    RRDDIM *rd = createTestDimension(rs);
    
    ml_dimension_new(rd);
    
    EXPECT_EQ(rd->ml_dimension, nullptr);
    
    freeDimension(rd);
    freeChart(rs);
    freeHost(rh);
}

// Tests for ml_dimension_delete()
TEST_F(MLPublicTest, ml_dimension_delete_frees_resources) {
    RRDHOST *rh = createTestHost();
    RRDSET *rs = createTestChart(rh);
    RRDDIM *rd = createTestDimension(rs);
    Cfg.enable_anomaly_detection = true;
    Cfg.num_worker_threads = 1;
    Cfg.workers.resize(1);
    Cfg.num_models_to_use = 1;
    
    ml_host_new(rh);
    ml_chart_new(rs);
    ml_dimension_new(rd);
    ml_dimension_delete(rd);
    
    EXPECT_EQ(rd->ml_dimension, nullptr);
    
    ml_chart_delete(rs);
    ml_host_delete(rh);
    freeDimension(rd);
    freeChart(rs);
    freeHost(rh);
}

TEST_F(MLPublicTest, ml_dimension_delete_handles_null_dimension) {
    RRDHOST *rh = createTestHost();
    RRDSET *rs = createTestChart(rh);
    RRDDIM *rd = createTestDimension(rs);
    
    ml_dimension_delete(rd);
    
    EXPECT_EQ(rd->ml_dimension, nullptr);
    
    freeDimension(rd);
    freeChart(rs);
    freeHost(rh);
}

// Tests for ml_dimension_received_anomaly()
TEST_F(MLPublicTest, ml_dimension_received_anomaly_handles_null_dimension) {
    RRDHOST *rh = createTestHost();
    RRDSET *rs = createTestChart(rh);
    RRDDIM *rd = createTestDimension(rs);
    
    ml_dimension_received_anomaly(rd, true);
    
    freeDimension(rd);
    freeChart(rs);
    freeHost(rh);
}

TEST_F(MLPublicTest, ml_dimension_received_anomaly_handles_not_running_host) {
    RRDHOST *rh = createTestHost();
    RRDSET *rs = createTestChart(rh);
    RRDDIM *rd = createTestDimension(rs);
    Cfg.enable_anomaly_detection = true;
    Cfg.num_worker_threads = 1;
    Cfg.workers.resize(1);
    Cfg.num_models_to_use = 1;
    
    ml_host_new(rh);
    ml_chart_new(rs);
    ml_dimension_new(rd);
    ml_dimension_received_anomaly(rd, true);
    
    ml_dimension_delete(rd);
    ml_chart_delete(rs);
    ml_host_