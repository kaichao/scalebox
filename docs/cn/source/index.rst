.. Scalebox documentation master file, created by
   sphinx-quickstart on Sun Jan 21 09:48:07 2024.
   You can adapt this file completely to your liking, but it should at least
   contain the root `toctree` directive.

Scalebox 中文文档
====================================

本中文站面向用户、应用开发者与系统管理员；AI agent 工作流与机器可读参考见仓库 `docs/agents/`。

根据您的使用场景，选择以下学习路径：

- **应用开发者**：:doc:`入门教程 <started/1_introduction>` → :doc:`核心概念 <user/2_core_concepts>` → :doc:`运行与管理应用 <user/3_running_apps>` → :doc:`应用设计 <developer/3_app_design>` → :doc:`模块开发 <developer/2_module_development>`
- **部署与运维**：:doc:`安装部署 <user/1_installation>` → :doc:`运维管理 <user/5_operations>` → :doc:`FAQ <faq>`
- **平台贡献者 / 深入原理**：:doc:`编程模型 <developer/1_programming_model>` → :doc:`主路由与应用状态 <developer/4_main_router_status>` → :doc:`VTask 虚拟任务 <developer/9_vtask>` → :doc:`节点本地计算 <developer/5_node_local_compute>` → :doc:`跨集群架构 <developer/10_cross_cluster>` → :doc:`安全框架 <developer/11_security>` → :doc:`高级编程特性 <developer/7_advanced_features>`

**开发者与 AI 参考（仓库直读）**

中文站面向人（用户、应用开发者、系统管理员）；AI agent 工作流与机器可读参考在仓库中直读：

- `AGENTS.md <https://github.com/kaichao/scalebox/blob/dev/AGENTS.md>`_
- `docs/agents/ <https://github.com/kaichao/scalebox/tree/dev/docs/agents>`_
- `features/ <https://github.com/kaichao/scalebox/tree/dev/features>`_
- `docs/reference/ <https://github.com/kaichao/scalebox/tree/dev/docs/reference>`_

.. toctree::
   :maxdepth: 2
   :caption: 一、入门教程
   :hidden:

   started/1_introduction
   started/2_quick_start
   started/3_first_app

.. toctree::
   :maxdepth: 2
   :caption: 二、应用开发
   :hidden:

   user/2_core_concepts
   user/3_running_apps
   developer/3_app_design
   developer/2_module_development
   user/4_example_apps
   user/6_standard_modules
   user/7_webui
   user/8_vscode

.. toctree::
   :maxdepth: 2
   :caption: 三、部署与运维
   :hidden:

   user/1_installation
   user/5_operations

.. toctree::
   :maxdepth: 2
   :caption: 四、设计与原理
   :hidden:

   developer/1_programming_model
   developer/4_main_router_status
   developer/9_vtask
   developer/5_node_local_compute
   developer/10_cross_cluster
   developer/11_security
   developer/7_advanced_features

.. toctree::
   :maxdepth: 2
   :caption: 五、参考
   :hidden:

   appendix/1_tech_specifications
   appendix/2_parameter_reference
   appendix/3_commandline_tools
   appendix/6_exit_code_spec

.. toctree::
   :maxdepth: 2
   :caption: 六、附录与其他
   :hidden:

   appendix/4_shell_programming
   appendix/7_software_install
   faq
   about_us/contributors
   about_us/contact_us

Indices and tables
==================

* :ref:`genindex`
* :ref:`modindex`
* :ref:`search`
