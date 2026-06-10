-- MySQL dump 10.13  Distrib 8.0.45, for Linux (x86_64)
--
-- Host: localhost    Database: library
-- ------------------------------------------------------
-- Server version	8.0.45-0ubuntu0.22.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `adjust_books`
--

DROP TABLE IF EXISTS `adjust_books`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `adjust_books` (
  `adjust_id` int NOT NULL AUTO_INCREMENT,
  `adjust_date` varchar(50) DEFAULT '',
  `adjust_title` varchar(500) DEFAULT '',
  `adjust_isbn` varchar(100) DEFAULT '',
  `adjust_content` text,
  PRIMARY KEY (`adjust_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `adjust_books`
--

LOCK TABLES `adjust_books` WRITE;
/*!40000 ALTER TABLE `adjust_books` DISABLE KEYS */;
INSERT INTO `adjust_books` VALUES (1,'2024-10-29','《活着》可借数量增加','1','《活着》作者：余华，可借数量增加3本'),(2,'2024-10-29','《狂人日记》存储位置调整','3','《狂人日记》存储位置调整到图书馆3楼7-101'),(3,'2024-10-29','《活着》书库数量增加','1','11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111'),(4,'2024-10-31','新书上架','10','村上春树的《挪威的森林》现以上架');
/*!40000 ALTER TABLE `adjust_books` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `admin`
--

DROP TABLE IF EXISTS `admin`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admin` (
  `admin_id` varchar(100) NOT NULL,
  `admin_password` varchar(255) NOT NULL DEFAULT '',
  `admin_role` varchar(50) DEFAULT '',
  PRIMARY KEY (`admin_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `admin`
--

LOCK TABLES `admin` WRITE;
/*!40000 ALTER TABLE `admin` DISABLE KEYS */;
INSERT INTO `admin` VALUES ('a','1','admin'),('b','1','Admin');
/*!40000 ALTER TABLE `admin` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `all_books`
--

DROP TABLE IF EXISTS `all_books`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `all_books` (
  `title` varchar(500) NOT NULL DEFAULT '',
  `author` varchar(255) DEFAULT '',
  `book_type` varchar(100) DEFAULT '',
  `press` varchar(255) DEFAULT '',
  `press_date` varchar(50) DEFAULT '',
  `isbn` varchar(100) NOT NULL,
  `cover` varchar(500) DEFAULT '',
  `intro` text,
  `price` decimal(10,2) DEFAULT '0.00',
  `amount` int DEFAULT '0',
  `lend_amount` int DEFAULT '0',
  `cur_lend_amount` int DEFAULT '0',
  `rec_state` int DEFAULT '0',
  PRIMARY KEY (`isbn`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `all_books`
--

LOCK TABLES `all_books` WRITE;
/*!40000 ALTER TABLE `all_books` DISABLE KEYS */;
INSERT INTO `all_books` VALUES ('活着','余华','文学类','人民出版社','2023-10-09','1','images/image1.jpg','这是一本好书',20.00,5,5,5,1),('挪威的森林','村上春树','科幻类','清华大学出版社','2023-10-18','10','images/image10.jpg','森林风光',25.00,3,3,3,1),('三体','刘慈欣','科幻类','清华大学出版社','2023','11','images/book-11.jpg','文化大革命如火如荼进行的同时，军方探寻外星文明的绝秘计划取得了突破性进展。但在按下发射键的那一刻，彻底改变了人类的命运。',29.00,5,3,2,1),('围城','钱钟书','文学类','人民出版社','2023','12','images/book-12.jpg','以抗战初期为背景，讲述了方鸿渐留学回国后，在动荡不安的社会中遭遇到的各种人生困境和爱情纠葛。',24.00,5,3,2,0),('平凡的世界','路遥','文学类','北京大学出版社','2022','13','images/book-13.jpg','以中国70年代中期到80年代中期十年间为背景，以孙少安和孙少平两兄弟为中心，展现了当时社会各阶层众多普通人的形象。',26.00,5,3,3,1),('红楼梦','曹雪芹','文学类','人民出版社','2024','14','images/book-14.jpg','以贾宝玉、林黛玉、薛宝钗的爱情婚姻悲剧为主线，描写了以贾家为代表的四大家族的兴衰历程。',28.00,5,3,3,0),('西游记','吴承恩','科幻类','人民邮电出版社','2023','15','images/book-15.jpg','孙悟空、猪八戒、沙僧三人保护唐僧西行取经，沿途历经八十一难，一路降妖伏魔，化险为夷，最后到达西天取得真经。',25.00,5,3,3,0),('白鹿原','陈忠实','历史类','人民邮电出版社','2021','16','images/book-16.jpg','以陕西关中平原上的白鹿村为背景，讲述了白姓和鹿姓两大家族祖孙三代的故事，展现了中国农村的广阔生活画面。',25.00,5,3,3,1),('射雕英雄传','金庸','历史类','清华大学出版社','2022','17','images/book-17.jpg','以宁宗庆元年至成吉思汗逝世为历史背景，讲述了傻小子郭靖背负家恨国仇闯入江湖，经历无数磨难成为一代大侠的故事。',28.00,5,3,3,0),('明朝那些事儿','当年明月','历史类','北京大学出版社','2024','18','images/book-18.jpg','以史料为基础，以年代和具体人物为主线，对明朝十七帝和其他王公权贵和小人物的命运进行全景展示。',28.00,5,3,3,1),('提问的艺术','[美] 特里.费德姆','文学类','人民邮电出版社','2024-10-06','2','images/image2.jpg','怎么提一个好问题问,看这本书就够了',20.00,4,4,4,0),('狂人日记','鲁迅','文学类','人民出版社','2024-10-01','3','images/image3.jpg','好!',20.00,5,5,5,1),('钢铁是怎样炼成的','(苏)尼·奥斯特洛夫斯基','文学类','人民出版社','2024-10-03','4','images/image4.jpg','钢铁般的意志',20.00,3,3,3,1),('骆驼祥子','老舍','文学类','清华大学出版社','2024-10-02','5','images/image5.jpg','祥子的一生',20.00,5,5,5,1),('海边的卡夫卡','村上春树','艺术类','人民出版社','2024-10-04','6','images/image6.jpg','卡卡',20.00,4,4,4,1),('百年孤独','加西亚·马尔克斯','历史类','北京大学出版社','2022-10-30','7','images/image7.jpg','孤独的人',22.00,5,5,5,1),('追风筝的人','卡勒德·胡赛尼','文学类','人民出版社','2023','8','images/image8.jpg','12岁的阿富汗富家少爷阿米尔与仆人哈桑的故事。一个关于人性、背叛与救赎的感人故事。',20.00,5,3,3,1),('动物农场','乔治·奥威尔','文学类','人民邮电出版社','2022','9','images/image9.jpg','一个农场的动物不堪人类压迫，在猪的带领下起来反抗，建立了平等的动物社会。然而动物领袖最终却成为比人类更独裁的统治者。',22.00,5,3,3,0);
/*!40000 ALTER TABLE `all_books` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `book_contents`
--

DROP TABLE IF EXISTS `book_contents`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `book_contents` (
  `isbn` varchar(100) NOT NULL,
  `content` longtext,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`isbn`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `book_contents`
--

LOCK TABLES `book_contents` WRITE;
/*!40000 ALTER TABLE `book_contents` DISABLE KEYS */;
INSERT INTO `book_contents` VALUES ('1','第一章\n\n阳光透过窗帘的缝隙洒在书桌上，我翻开这本泛黄的书页，试图在文字间寻找那些逝去的时光。\n\n故事要从那个秋天说起。那时候我还年轻，对世界充满好奇。每天骑着自行车穿过城市的街道，看着路边的梧桐树叶一片片变黄、飘落。\n\n人生就像一场旅行，不必在乎目的地，在乎的是沿途的风景以及看风景的心情。\n（更多内容等待管理员上传...）','2026-06-02 14:57:56'),('11','第一章 疯狂年代\n\n\n　　“文革”中经历过劫难的人大都记得，那段岁月有一种说法叫做“清理阶级队伍”。\n\n　　事情的起因是，军方探寻外星文明的绝秘计划“红岸工程”取得了突破性进展。但在按下发射键的那一刻，历经劫难的叶文洁没有意识到，她彻底改变了人类的命运。\n\n　　地球文明向宇宙发出的第一声啼鸣，以太阳为中心，以光速向宇宙深处飞驰……\n\n　　四光年外，“三体文明”正苦苦挣扎——三颗无规则运行的太阳主导下的百余次毁灭与重生逼迫他们逃离自己的家园。而恰在此时，他们接收到了地球发来的信息。\n\n　　对人性绝望的叶文洁向三体人暴露了地球的坐标，一场即将到来的文明浩劫，从现在开始。\n\n（未完待续，更多内容等待管理员上传）','2026-06-02 15:57:24'),('12','第一章\n\n\n　　红海早过了，所以船到香港特别慢。\n\n　　方鸿渐到香港，一个叫苏文纨的女人来接他。她曾是他的大学同学，现在在香港做官太太。\n\n　　方鸿渐在欧洲留学四年，换了三个大学，最后从一个叫“克莱登大学”的野鸡学校买了张假博士文凭回国。\n\n　　在船上，一个姓鲍的性感女人引诱了他。到了香港，苏文纨又对他示好。方鸿渐周旋其间，疲于应付。\n\n　　回到上海，方鸿渐住进已故未婚妻的父家，在岳父的点金银行任职。他又认识了苏文纨的表妹唐晓芙，一见倾心。\n\n　　然而，在这座围城中，有人想进去，有人想出来。方鸿渐的爱情、事业、婚姻，都将在这座围城中经受考验……\n\n（未完待续，更多内容等待管理员上传）','2026-06-02 15:57:24'),('13','第一章\n\n\n　　一九七五年二三月间，一个平平常常的日子，细蒙蒙的雨丝夹着一星半点的雪花，正纷纷淋淋地向大地飘洒着。\n\n　　时令已快到惊蛰，雪当然再不会存留，往往还没等落地，就已经消失得无踪无影了。黄土高原严寒而漫长的冬天看来就要过去，但那真正温暖的春天还远远地没有到来。\n\n　　在这样雨雪交加的日子里，如果没有什么紧要事，人们宁愿一整天足不出户。\n\n　　县城的大街上，只有寥寥几个匆匆赶路的人。孙少平来到学校，他是全县最好的高中生之一，却因为家境贫寒，每天只能吃两个黑馍馍……\n\n　　但正是这个平凡的青年，用他坚韧不拔的意志，书写了一段不平凡的人生。\n\n（未完待续，更多内容等待管理员上传）','2026-06-02 15:57:24'),('14','第一回 甄士隐梦幻识通灵 贾雨村风尘怀闺秀\n\n\n　　此开卷第一回也。作者自云：因曾历过一番梦幻之后，故将真事隐去，而借“通灵”之说，撰此《石头记》一书也。\n\n　　列位看官：你道此书从何而来？说起根由虽近荒唐，细按则深有趣味。\n\n　　原来女娲氏炼石补天之时，于大荒山无稽崖练成高经十二丈、方经二十四丈顽石三万六千五百零一块。娲皇氏只用了三万六千五百块，只单单剩了一块未用，便弃在此山青埂峰下。\n\n　　谁知此石自经煅炼之后，灵性已通，因见众石俱得补天，独自己无材不堪入选，遂自怨自叹，日夜悲号惭愧。\n\n　　一日，一僧一道来到这大荒山，说起红尘中荣华富贵，此石动了凡心，想要到人间去享一享……\n\n（未完待续，更多内容等待管理员上传）','2026-06-02 15:57:24'),('15','第一回 灵根育孕源流出 心性修持大道生\n\n\n　　诗曰：\n　　混沌未分天地乱，茫茫渺渺无人见。\n　　自从盘古破鸿蒙，开辟从兹清浊辨。\n　　覆载群生仰至仁，发明万物皆成善。\n　　欲知造化会元功，须看西游释厄传。\n\n　　海外有一国土，名曰傲来国。国近大海，海中有一座名山，唤为花果山。\n\n　　那座山正当顶上，有一块仙石。自开辟以来，每受天真地秀，日精月华，感之既久，遂有灵通之意。内育仙胞，一日迸裂，产一石卵，似圆球样大。因见风，化作一个石猴，五官俱备，四肢皆全。\n\n　　目运金光，射冲斗府。惊动高天上圣大慈仁者玉皇大天尊玄穹高上帝……\n\n（未完待续，更多内容等待管理员上传）','2026-06-02 15:57:24'),('16','第一章\n\n\n　　白嘉轩后来引以为豪壮的是一生里娶过七房女人。\n\n　　但那是从前的光景了。如今的西安府白鹿原上，白家和鹿家两姓人世代聚居于此。白嘉轩作为族长，恪守着祖辈传下的规矩，勤恳耕作，善待乡邻。\n\n　　时代在变，清朝覆灭，民国建立，革命浪潮席卷全国。白鹿原上的人们，也不可避免地被卷入了历史的洪流。\n\n　　白嘉轩的儿子白孝文和鹿子霖的儿子鹿兆鹏，走上了截然不同的人生道路。一个固守传统，一个投身革命。\n\n　　在半个世纪的历史变迁中，白鹿两家几代人的恩怨纠葛，在这片古老的土地上缓缓展开……\n\n（未完待续，更多内容等待管理员上传）','2026-06-02 15:57:24'),('17','第一回 风雪惊变\n\n\n　　钱塘江浩浩江水，日日夜夜无穷无休地奔向大海。\n\n　　越女采莲秋水畔，窄袖轻罗，暗露双金钏。照影摘花花似面，芳心只共丝争乱。\n\n　　南宋宁宗庆元年间，一个叫郭啸天的好汉，在临安城外牛家村与义弟杨铁心相会。两人一个是梁山好汉郭盛的后代，一个是抗金名将杨再兴的后人。\n\n　　不想这一聚，却引来了一场灭顶之灾。\n\n　　郭啸天惨死，其妻逃到大漠，生下遗腹子郭靖。杨铁心生死不明，其妻被掳，生下杨康。\n\n　　十八年后，郭靖已长成一个憨厚朴实的草原少年。他虽天资愚钝，但心地纯良，机缘巧合之下，拜江南七怪为师，又得全真教马钰道长传授内功。\n\n　　从此，这个傻小子将一步步走进江湖，开启一段波澜壮阔的英雄传奇……\n\n（未完待续，更多内容等待管理员上传）','2026-06-02 15:57:24'),('18','引子\n\n\n　　很多人问，为什么要读历史？\n\n　　其实这个问题很简单，因为历史离我们很近。\n\n　　比如明朝那些事儿，离我们今天不过几百年。几百年前的人，和我们一样，要吃饭，要睡觉，要过日子。他们也会高兴，也会悲伤，也会愤怒，也会无奈。\n\n　　不同的是，他们生活在那个时代，而我们生活在这个时代。\n\n　　真正的历史，不是枯燥的年代和人名，而是鲜活的人和他们做的事。\n\n　　所以，这本书就是要告诉你——\n\n　　历史本身很精彩，历史可以写得很好看。\n\n　　让我们从明朝开国皇帝朱元璋说起。他出身贫农，当过和尚，讨过饭，最后却成为开创一个王朝的伟大人物。\n\n　　他的人生，本身就是一部传奇……\n\n（未完待续，更多内容等待管理员上传）','2026-06-02 15:57:24'),('2','提问是一门艺术。\n\n一个好的问题，往往比答案更重要。因为问题是思考的起点，是探索未知的动力。\n\n从前有一个年轻人，他总是有很多问题。他去问智者：\"如何才能获得智慧？\"\n智者说：\"继续提问。\"\n年轻人又问：\"那我该问什么样的问题？\"\n智者说：\"问那些让你睡不着觉的问题。\"\n—— 这就是提问的艺术。\n（更多内容等待管理员上传...）','2026-06-02 14:57:56'),('3','今天晚上，很好的月光。\n\n我不见他，已是三十多年；今天见了，精神分外爽快。才知道以前的三十多年，全是发昏；然而须十分小心。不然，那赵家的狗，何以看我两眼呢？\n\n我怕得有理。\n（更多内容等待管理员上传...）','2026-06-02 14:57:56');
/*!40000 ALTER TABLE `book_contents` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cur_lend_records`
--

DROP TABLE IF EXISTS `cur_lend_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cur_lend_records` (
  `lend_id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL DEFAULT '',
  `title` varchar(500) NOT NULL DEFAULT '',
  `isbn` varchar(100) NOT NULL DEFAULT '',
  `lend_date` varchar(50) DEFAULT '',
  `exp_return_date` varchar(50) DEFAULT '',
  PRIMARY KEY (`lend_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cur_lend_records`
--

LOCK TABLES `cur_lend_records` WRITE;
/*!40000 ALTER TABLE `cur_lend_records` DISABLE KEYS */;
INSERT INTO `cur_lend_records` VALUES (2,'1','三体','11','2026-06-02','2026-06-03'),(3,'1','围城','12','2026-06-02','2026-06-03');
/*!40000 ALTER TABLE `cur_lend_records` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `lend_records`
--

DROP TABLE IF EXISTS `lend_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `lend_records` (
  `lend_id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL DEFAULT '',
  `title` varchar(500) NOT NULL DEFAULT '',
  `isbn` varchar(100) NOT NULL DEFAULT '',
  `lend_date` varchar(50) DEFAULT '',
  `exp_return_date` varchar(50) DEFAULT '',
  PRIMARY KEY (`lend_id`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `lend_records`
--

LOCK TABLES `lend_records` WRITE;
/*!40000 ALTER TABLE `lend_records` DISABLE KEYS */;
INSERT INTO `lend_records` VALUES (1,'1','活着','1','2024-10-22','2024-10-26'),(2,'1','提问的艺术','2','2024-10-22','2024-10-26'),(3,'1','狂人日记','3','2024-10-22','2024-10-26'),(4,'1','活着','1','2024-10-22','2024-10-24'),(5,'1','钢铁是怎样炼成的','4','2024-10-25','2024-10-25'),(6,'2','狂人日记','3','2024-10-29','2024-10-29'),(7,'1','活着','1','2024-10-29','2024-10-29'),(8,'1','海边的卡夫卡','6','2024-10-29','2024-10-29'),(9,'2','活着','1','2024-10-30','2024-10-30'),(10,'2','提问的艺术','2','2024-10-30','2024-10-30'),(11,'1','活着','1','2024-10-31','2024-10-31'),(12,'1','提问的艺术','2','2024-10-31','2024-10-31'),(13,'1','挪威的森林','10','2024-10-31','2024-10-31'),(14,'1','百年孤独','7','2024-10-31','2024-10-31'),(15,'1','挪威的森林','10','2024-10-31','2024-11-01'),(16,'1','提问的艺术','2','2024-10-31','2024-11-01'),(17,'1','活着','1','2024-10-31','2024-11-01'),(18,'1','活着','1','2024-10-31','2024-10-31'),(19,'2','活着','1','2024-11-05','2024-11-05'),(20,'1','活着','1','2024-11-05','2024-11-05'),(21,'1','活着','1','2025-07-04','2025-07-16'),(22,'1','活着','1','2025-12-05','2025-12-12'),(23,'1','活着','1','2026-03-17','2026-03-21'),(24,'1','钢铁是怎样炼成的','4','2026-04-01','2026-04-04'),(25,'1','活着','1','2026-04-07','2026-04-07'),(26,'1','活着','1','2026-06-02','2026-06-16'),(27,'1','三体','11','2026-06-02','2026-06-03'),(28,'1','围城','12','2026-06-02','2026-06-03');
/*!40000 ALTER TABLE `lend_records` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `library_summary`
--

DROP TABLE IF EXISTS `library_summary`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `library_summary` (
  `total_books_amount` int DEFAULT '0',
  `total_lend_amount` int DEFAULT '0',
  `total_return_amount` int DEFAULT '0',
  `total_users_amount` int DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `library_summary`
--

LOCK TABLES `library_summary` WRITE;
/*!40000 ALTER TABLE `library_summary` DISABLE KEYS */;
INSERT INTO `library_summary` VALUES (34,32,2,4);
/*!40000 ALTER TABLE `library_summary` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notices`
--

DROP TABLE IF EXISTS `notices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `notices` (
  `notice_id` int NOT NULL AUTO_INCREMENT,
  `notice_date` varchar(50) DEFAULT '',
  `notice_title` varchar(500) DEFAULT '',
  `notice` text,
  PRIMARY KEY (`notice_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notices`
--

LOCK TABLES `notices` WRITE;
/*!40000 ALTER TABLE `notices` DISABLE KEYS */;
INSERT INTO `notices` VALUES (3,'2024-10-31','测试公告','测试');
/*!40000 ALTER TABLE `notices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `recommend_books`
--

DROP TABLE IF EXISTS `recommend_books`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `recommend_books` (
  `isbn` varchar(100) NOT NULL,
  `title` varchar(500) NOT NULL DEFAULT '',
  `author` varchar(255) DEFAULT '',
  `rec_type` varchar(100) DEFAULT '',
  `cover` varchar(500) DEFAULT '',
  `cur_lend_amount` int DEFAULT '0',
  PRIMARY KEY (`isbn`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `recommend_books`
--

LOCK TABLES `recommend_books` WRITE;
/*!40000 ALTER TABLE `recommend_books` DISABLE KEYS */;
INSERT INTO `recommend_books` VALUES ('1','活着','余华','新书','images/image1.jpg',5),('10','挪威的森林','村上春树','新书','images/image10.jpg',3),('11','三体','刘慈欣','科幻经典','images/book-11.jpg',2),('13','平凡的世界','路遥','经典文学','images/book-13.jpg',3),('16','白鹿原','陈忠实','经典文学','images/book-16.jpg',3),('18','明朝那些事儿','当年明月','经典历史','images/book-18.jpg',3),('3','狂人日记','鲁迅','新书','images/image3.jpg',5),('4','钢铁是怎样炼成的','(苏)尼·奥斯特洛夫斯基','新书','images/image4.jpg',3),('5','骆驼祥子','老舍','新书','images/image5.jpg',5),('6','海边的卡夫卡','村上春树','新书','images/image6.jpg',4),('7','百年孤独','加西亚·马尔克斯','新书','images/image7.jpg',5),('8','追风筝的人','卡勒德·胡赛尼','推荐好书','images/image8.jpg',3);
/*!40000 ALTER TABLE `recommend_books` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `refresh_tokens`
--

DROP TABLE IF EXISTS `refresh_tokens`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `refresh_tokens` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `token_hash` varchar(255) NOT NULL,
  `expires_at` datetime NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `refresh_tokens`
--

LOCK TABLES `refresh_tokens` WRITE;
/*!40000 ALTER TABLE `refresh_tokens` DISABLE KEYS */;
INSERT INTO `refresh_tokens` VALUES (24,'1','$2a$10$PterflKoZyhRJMmawvTqxuE/I7khaI5x74fj80ZVEssKElicHEqzu','2026-06-09 15:06:44','2026-06-02 15:06:44'),(25,'1','$2a$10$CRdRwSCTDfLVtAw7nAK0euDlGz838p83sb39ENFkbnABofqg/WHqy','2026-06-09 15:57:33','2026-06-02 15:57:33'),(26,'1','$2a$10$L8WLJ7Y9E91tH3uJywWHheBwce63kDyyn3LkG6.EQxIoJOFSv7ZF2','2026-06-09 15:58:20','2026-06-02 15:58:20'),(27,'1','$2a$10$ewMhnwnP27OrFWVAEKseGub3wGt/psbn4COZjS9h/ZjY.kEwNTrD2','2026-06-09 16:05:46','2026-06-02 16:05:46'),(28,'1','$2a$10$Yjy8SnLUyBvsaTbwwsoFEuNTfBIgQjuLgunpL7.SsaLhPvdKyudp.','2026-06-09 16:13:48','2026-06-02 16:13:48');
/*!40000 ALTER TABLE `refresh_tokens` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `replay_opinions`
--

DROP TABLE IF EXISTS `replay_opinions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `replay_opinions` (
  `replay_id` int NOT NULL AUTO_INCREMENT,
  `replay_name` varchar(255) DEFAULT '',
  `replay_date` varchar(50) DEFAULT '',
  `replay_idea` text,
  `replay_user` varchar(255) DEFAULT '',
  PRIMARY KEY (`replay_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `replay_opinions`
--

LOCK TABLES `replay_opinions` WRITE;
/*!40000 ALTER TABLE `replay_opinions` DISABLE KEYS */;
INSERT INTO `replay_opinions` VALUES (1,'智慧图书馆','2024-10-30','1','1'),(2,'智慧图书馆','2024-10-31','1','1'),(3,'智慧图书馆','2026-03-25','good','1');
/*!40000 ALTER TABLE `replay_opinions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `return_records`
--

DROP TABLE IF EXISTS `return_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `return_records` (
  `return_id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL DEFAULT '',
  `title` varchar(500) NOT NULL DEFAULT '',
  `isbn` varchar(100) NOT NULL DEFAULT '',
  `lend_date` varchar(50) DEFAULT '',
  `exp_return_date` varchar(50) DEFAULT '',
  `return_date` varchar(50) DEFAULT '',
  `late_fee` decimal(10,2) DEFAULT '0.00',
  PRIMARY KEY (`return_id`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `return_records`
--

LOCK TABLES `return_records` WRITE;
/*!40000 ALTER TABLE `return_records` DISABLE KEYS */;
INSERT INTO `return_records` VALUES (1,'1','狂人日记','3','2024-10-22','2024-10-26','2024-10-22',0.00),(2,'1','活着','1','2024-10-22','2024-10-26','2024-10-22',0.00),(3,'1','活着','1','2024-10-22','2024-10-24','2024-10-25',0.40),(4,'1','钢铁是怎样炼成的','4','2024-10-25','2024-10-25','2024-10-28',1.20),(5,'1','提问的艺术','2','2024-10-22','2024-10-26','2024-10-28',0.80),(6,'1','活着','1','2024-10-29','2024-10-29','2024-10-31',0.80),(7,'1','海边的卡夫卡','6','2024-10-29','2024-10-29','2024-10-31',0.80),(8,'1','活着','1','2024-10-31','2024-10-31','2024-10-31',0.00),(9,'1','提问的艺术','2','2024-10-31','2024-10-31','2024-10-31',0.00),(10,'2','活着','1','2024-10-30','2024-10-30','2024-10-31',0.40),(11,'2','狂人日记','3','2024-10-29','2024-10-29','2024-10-31',0.80),(12,'2','提问的艺术','2','2024-10-30','2024-10-30','2024-10-31',0.40),(13,'1','挪威的森林','10','2024-10-31','2024-10-31','2024-10-31',0.00),(14,'1','活着','1','2024-10-31','2024-11-01','2024-10-31',0.00),(15,'1','百年孤独','7','2024-10-31','2024-10-31','2024-10-31',0.00),(16,'1','挪威的森林','10','2024-10-31','2024-11-01','2024-10-31',0.00),(17,'1','提问的艺术','2','2024-10-31','2024-11-01','2024-10-31',0.00),(18,'1','活着','1','2024-10-31','2024-10-31','2024-10-31',0.00),(19,'2','活着','1','2024-11-05','2024-11-05','2024-11-05',0.00),(20,'1','活着','1','2024-11-05','2024-11-05','2024-11-05',0.00),(21,'1','活着','1','2025-07-04','2025-07-16','2025-12-05',56.80),(22,'1','活着','1','2025-12-05','2025-12-12','2026-03-17',38.00),(23,'1','活着','1','2026-03-17','2026-03-21','2026-03-19',0.00),(24,'1','钢铁是怎样炼成的','4','2026-04-01','2026-04-04','2026-04-07',1.20),(25,'1','活着','1','2026-04-07','2026-04-07','2026-04-16',3.60),(26,'1','活着','1','2026-06-02','2026-06-16','2026-06-02',0.00);
/*!40000 ALTER TABLE `return_records` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `session_state`
--

DROP TABLE IF EXISTS `session_state`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `session_state` (
  `session_name` varchar(255) DEFAULT NULL,
  `session_id` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `session_state`
--

LOCK TABLES `session_state` WRITE;
/*!40000 ALTER TABLE `session_state` DISABLE KEYS */;
INSERT INTO `session_state` VALUES ('a','2068ee0a-0f6e-4a6f-a979-a4186afbcbc4'),('1','45c04180-e92a-4480-a2af-0c59e4758bf6');
/*!40000 ALTER TABLE `session_state` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_opinions`
--

DROP TABLE IF EXISTS `user_opinions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_opinions` (
  `opinion_id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT '',
  `phone` varchar(100) DEFAULT '',
  `email` varchar(255) DEFAULT '',
  `idea` text,
  PRIMARY KEY (`opinion_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_opinions`
--

LOCK TABLES `user_opinions` WRITE;
/*!40000 ALTER TABLE `user_opinions` DISABLE KEYS */;
INSERT INTO `user_opinions` VALUES (1,'12','123','333','111'),(2,'1','1','1','1'),(3,'11111111111111111111111111','222222222222222222','33333333333333333','555555555555555555555555'),(4,'1','2','2','4');
/*!40000 ALTER TABLE `user_opinions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `name` varchar(255) NOT NULL DEFAULT '',
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL DEFAULT '',
  `user_cur_lend_amount` int DEFAULT '0',
  `user_his_lend_amount` int DEFAULT '0',
  `birthday` varchar(50) DEFAULT '',
  `age` int DEFAULT '0',
  `photo` varchar(500) DEFAULT '',
  `email` varchar(255) DEFAULT '',
  `email_verified` tinyint(1) DEFAULT '0',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES ('嘟嘟','1','$2a$10$ldXy2C211vxTYWLA8cTohOugEx7EByGmVGCDub5JrD4oo7Xfx3/jy',2,24,'2024-10-25',1,'images/image-1.jpg','2856791147@qq.com',1,'2026-06-02 13:00:28'),('奕奕','2','$2a$10$h4nQOUGT626SpnC8ysT8deYMuBDwlQnuM/ecg5MtbFadzHTUR/zrC',0,4,'2024-10-25',1,'images/image-2.jpg','',1,'2026-06-02 13:00:28'),('3368820338@qq.com','3368820338@qq.com','$2a$10$jk03GYqfCxEBD9byE1CaLeUYCQO5Zg5QnnLqPM9gq4NdLXx5fKfBi',0,0,'',0,'','3368820338@qq.com',1,'2026-06-02 14:22:42');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `verification_codes`
--

DROP TABLE IF EXISTS `verification_codes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `verification_codes` (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  `code` varchar(6) NOT NULL,
  `expires_at` datetime NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `verification_codes`
--

LOCK TABLES `verification_codes` WRITE;
/*!40000 ALTER TABLE `verification_codes` DISABLE KEYS */;
INSERT INTO `verification_codes` VALUES (2,'test@test.com','588633','2026-06-02 13:56:18','2026-06-02 13:46:18'),(3,'test-new-user@qq.com','292906','2026-06-02 14:06:20','2026-06-02 13:56:20'),(4,'newuser-test@qq.com','821635','2026-06-02 14:07:30','2026-06-02 13:57:30'),(13,'unregistered@test.com','695483','2026-06-02 14:18:46','2026-06-02 14:08:46'),(16,'ratelimit@test.com','640007','2026-06-02 14:20:02','2026-06-02 14:10:02'),(18,'2856791147@qq.com','878064','2026-06-02 14:31:44','2026-06-02 14:21:44');
/*!40000 ALTER TABLE `verification_codes` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-06-02 16:26:15
