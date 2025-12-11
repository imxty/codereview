import logging
import logging.config
import os
import yaml
from datetime import datetime, timedelta
from models.database import db
from models.log_models import SystemLog, LogConfig
from utils.log_handlers import DatabaseHandler, ContextFilter, StructuredFormatter, JsonFormatter
import json

class LogService:
    def __init__(self):
        self.context_filter = ContextFilter()
        self.loggers = {}
        self.default_config = self._get_default_config()
        self.initialized = False
    
    def _get_default_config(self):
        """获取默认日志配置"""
        return {
            'version': 1,
            'disable_existing_loggers': False,
            'formatters': {
                'standard': {
                    'format': '%(asctime)s - %(name)s - %(levelname)s - %(module)s:%(funcName)s:%(lineno)d - %(message)s'
                },
                'json': {
                    '()': JsonFormatter
                },
                'simple': {
                    'format': '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
                }
            },
            'filters': {
                'context_filter': {
                    '()': ContextFilter
                }
            },
            'handlers': {
                'console': {
                    'class': 'logging.StreamHandler',
                    'level': 'INFO',
                    'formatter': 'standard',
                    'stream': 'ext://sys.stdout'
                },
                'file_handler': {
                    'class': 'logging.handlers.RotatingFileHandler',
                    'level': 'INFO',
                    'formatter': 'standard',
                    'filename': 'logs/test_platform.log',
                    'maxBytes': 10485760,  # 10MB
                    'backupCount': 5,
                    'encoding': 'utf8'
                },
                'error_file_handler': {
                    'class': 'logging.handlers.RotatingFileHandler',
                    'level': 'ERROR',
                    'formatter': 'standard',
                    'filename': 'logs/error.log',
                    'maxBytes': 10485760,
                    'backupCount': 5,
                    'encoding': 'utf8'
                }
            },
            'loggers': {
                '': {  # root logger
                    'level': 'INFO',
                    'handlers': ['console', 'file_handler', 'error_file_handler'],
                    'propagate': False
                },
                'test_platform': {
                    'level': 'DEBUG',
                    'handlers': ['console', 'file_handler'],
                    'propagate': False
                },
                'services': {
                    'level': 'DEBUG',
                    'handlers': ['console', 'file_handler'],
                    'propagate': False
                },
                'models': {
                    'level': 'INFO',
                    'handlers': ['file_handler'],
                    'propagate': False
                },
                'routes': {
                    'level': 'INFO',
                    'handlers': ['console', 'file_handler'],
                    'propagate': False
                }
            }
        }
    
    def initialize_logging(self):
        """初始化日志系统"""
        try:
            if self.initialized:
                return
                
            # 创建日志目录
            os.makedirs('logs', exist_ok=True)
            
            # 应用默认配置
            logging.config.dictConfig(self.default_config)
            
            # 为根日志记录器添加上下文过滤器
            root_logger = logging.getLogger()
            root_logger.addFilter(self.context_filter)
            
            # 添加数据库处理器到根日志记录器（只在应用上下文中）
            try:
                from flask import current_app
                if current_app:
                    db_handler = DatabaseHandler(level=logging.INFO)
                    db_handler.setFormatter(StructuredFormatter())
                    root_logger.addHandler(db_handler)
            except:
                pass  # 不在应用上下文中，跳过数据库处理器
            
            # 从数据库加载自定义配置
            self._load_database_configs()
            
            self.initialized = True
            logging.info("日志系统初始化完成")
            
        except Exception as e:
            print(f"日志系统初始化失败: {str(e)}")
            # 使用基础配置
            logging.basicConfig(
                level=logging.INFO,
                format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
            )
    
    def _load_database_configs(self):
        """从数据库加载日志配置"""
        try:
            from flask import current_app
            if not current_app:
                return
                
            with current_app.app_context():
                configs = LogConfig.query.all()
                for config in configs:
                    self._apply_logger_config(config)
        except Exception as e:
            logging.error(f"加载数据库日志配置失败: {str(e)}")
    
    def _apply_logger_config(self, config):
        """应用日志记录器配置"""
        try:
            logger = logging.getLogger(config.logger_name)
            logger.setLevel(getattr(logging, config.level))
            
            # 清除现有处理器
            for handler in logger.handlers[:]:
                logger.removeHandler(handler)
            
            # 添加控制台处理器
            if config.enable_console:
                console_handler = logging.StreamHandler()
                console_handler.setLevel(getattr(logging, config.level))
                console_handler.setFormatter(StructuredFormatter(config.format))
                logger.addHandler(console_handler)
            
            # 添加文件处理器
            if config.enable_file and config.file_path:
                # 确保目录存在
                os.makedirs(os.path.dirname(config.file_path), exist_ok=True)
                
                file_handler = logging.handlers.RotatingFileHandler(
                    filename=config.file_path,
                    maxBytes=config.max_file_size,
                    backupCount=config.backup_count,
                    encoding='utf8'
                )
                file_handler.setLevel(getattr(logging, config.level))
                file_handler.setFormatter(StructuredFormatter(config.format))
                logger.addHandler(file_handler)
            
            # 添加数据库处理器（只在应用上下文中）
            if config.enable_database:
                try:
                    from flask import current_app
                    if current_app:
                        db_handler = DatabaseHandler(level=getattr(logging, config.level))
                        db_handler.setFormatter(StructuredFormatter(config.format))
                        logger.addHandler(db_handler)
                except:
                    pass
            
            self.loggers[config.logger_name] = logger
            
        except Exception as e:
            logging.error(f"应用日志配置失败: {config.logger_name}, 错误: {str(e)}")
    
    def set_context(self, **kwargs):
        """设置日志上下文"""
        self.context_filter.set_context(**kwargs)
    
    def clear_context(self):
        """清除日志上下文"""
        self.context_filter.clear_context()
    
    def get_logger(self, name):
        """获取日志记录器"""
        return logging.getLogger(name)
    
    def log_execution_start(self, execution_id, suite_name, environment):
        """记录执行开始"""
        self.set_context(execution_id=execution_id)
        logging.info(f"开始执行测试套件: {suite_name} [环境: {environment}]")
    
    def log_execution_end(self, execution_id, suite_name, result):
        """记录执行结束"""
        self.set_context(execution_id=execution_id)
        logging.info(f"测试套件执行完成: {suite_name} [结果: {result}]")
        self.clear_context()
    
    def log_test_case_start(self, test_case_id, test_case_name, test_type):
        """记录测试用例开始"""
        self.set_context(test_case_id=test_case_id)
        logging.info(f"开始执行测试用例: {test_case_name} [类型: {test_type}]")
    
    def log_test_case_end(self, test_case_id, test_case_name, status, execution_time):
        """记录测试用例结束"""
        self.set_context(test_case_id=test_case_id)
        logging.info(f"测试用例执行完成: {test_case_name} [状态: {status}, 耗时: {execution_time:.2f}s]")
        self.clear_context()
    
    def log_api_request(self, test_case_id, method, url, headers, data):
        """记录API请求"""
        self.set_context(test_case_id=test_case_id)
        logging.debug(f"API请求: {method} {url}")
        if headers:
            logging.debug(f"请求头: {headers}")
        if data:
            logging.debug(f"请求数据: {data}")
    
    def log_api_response(self, test_case_id, status_code, response_time, response_data):
        """记录API响应"""
        self.set_context(test_case_id=test_case_id)
        logging.debug(f"API响应: 状态码={status_code}, 耗时={response_time:.2f}s")
        if response_data and len(str(response_data)) < 1000:  # 限制日志长度
            logging.debug(f"响应数据: {response_data}")
    
    def log_web_action(self, test_case_id, action, element, value=None):
        """记录Web操作"""
        self.set_context(test_case_id=test_case_id)
        if value:
            logging.debug(f"Web操作: {action} {element} = {value}")
        else:
            logging.debug(f"Web操作: {action} {element}")
    
    def log_app_action(self, test_case_id, action, element=None, value=None):
        """记录App操作"""
        self.set_context(test_case_id=test_case_id)
        if element and value:
            logging.debug(f"App操作: {action} {element} = {value}")
        elif element:
            logging.debug(f"App操作: {action} {element}")
        else:
            logging.debug(f"App操作: {action}")
    
    def log_error(self, test_case_id, error_message, exc_info=None):
        """记录错误信息"""
        self.set_context(test_case_id=test_case_id)
        logging.error(f"执行错误: {error_message}", exc_info=exc_info)
        self.clear_context()

# 全局日志服务实例
log_service = LogService()